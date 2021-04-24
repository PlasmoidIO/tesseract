package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"share/common/packet"
	"share/peer/protocol"
	"strings"
)

type CentralClient struct {
	Conn               net.Conn
	DataChannels       []chan []byte
	sendRequestHandler func(p *packet.SendPacket) bool
	Started            bool
	RegisteredUsername string
}

func catch(err error) {
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}

func NewClient() CentralClient {
	conn, err := net.Dial("tcp", "localhost:8080")
	catch(err)
	centralClient := CentralClient{
		Conn:         conn,
		DataChannels: []chan []byte{},
		Started:      false,
	}
	return centralClient
}

func (cl *CentralClient) Start() {
	cl.Started = true
	defer func() {
		catch(cl.Conn.Close())
	}()

	scanner := bufio.NewScanner(cl.Conn)
	for scanner.Scan() {
		go cl.handleData(scanner.Bytes())
	}
}

func (cl *CentralClient) WritePacket(p packet.Packet) {
	if !cl.Started {
		fmt.Println("Error: trying to write data before client connected.")
		return
	}

	data := p.String()
	fmt.Printf("Writing %s\n", data)
	_, err := fmt.Fprintln(cl.Conn, data)
	catch(err)
}

func (cl *CentralClient) CreateDataChannel() chan []byte {
	ch := make(chan []byte)
	cl.DataChannels = append(cl.DataChannels, ch)
	return ch
}

func (cl *CentralClient) RemoveDataChannel(c chan []byte) {
	for i, ch := range cl.DataChannels {
		if ch == c {
			last := len(cl.DataChannels) - 1
			cl.DataChannels[i], cl.DataChannels[last] = cl.DataChannels[last], cl.DataChannels[i]
			cl.DataChannels = cl.DataChannels[:last]
		}
	}
}

func (cl *CentralClient) handleData(buf []byte) {
	data := string(buf)
	packetType := packet.GetPacketType(data)
	if packetType == "FILE_SEND_REQUEST" {
		p := packet.ToSendPacket(data)
		handler := cl.sendRequestHandler
		if handler != nil {
			res := handler(p)
			if res {
				app := protocol.NewShareHandler()
				accepted := packet.NewAcceptPacket(p.Filename, p.Size, app.PeerHandler.GetPeerAddress())
				cl.WritePacket(&accepted)
				if err := app.Receive(p); err != nil {
					fmt.Println(err)
				}
			} else {
				rejected := packet.NewRejectPacket(p.Filename)
				cl.WritePacket(&rejected)
			}
		}
	}

	for _, c := range cl.DataChannels {
		c <- buf
	}
}

func (cl *CentralClient) sendFileRequest(req *packet.SendPacket) (*packet.AcceptPacket, bool) {
	cl.WritePacket(req)
	ch := cl.CreateDataChannel()
	for {
		res := string(<-ch)
		packetType := packet.GetPacketType(res)
		if packetType == packet.ACCEPT {
			cl.RemoveDataChannel(ch)
			acceptPacket := packet.ToAcceptPacket(res)
			return acceptPacket, true
		}
		if packetType == packet.REJECT {
			rejectPacket := packet.ToRejectPacket(res)
			if rejectPacket.Filename == req.Filename {
				cl.RemoveDataChannel(ch)
				return nil, false
			}
		}
	}
}

func (cl *CentralClient) SendFile(filename string, filesize int, target string) bool {
	app := protocol.NewShareHandler()
	req := packet.NewSendPacket(filename, filesize, target, app.PeerHandler.GetPeerAddress())
	accept, ok := cl.sendFileRequest(&req)
	if ok {
		catch(app.Send(accept))
	}
	return ok
}

func (cl *CentralClient) RegisterUsername(username string) error {
	registerPacket := packet.NewRegisterPacket(username)
	c := cl.CreateDataChannel()
	cl.WritePacket(&registerPacket)
	for {
		res := <-c
		data := strings.Split(string(res), " ")
		if len(data) < 2 {
			continue
		}
		if data[0] == registerPacket.PacketType {
			cl.RemoveDataChannel(c)
			switch data[1] {
			case "USER_REGISTER_SUCCESS":
				cl.RegisteredUsername = username
				return nil
			case "USER_REGISTER_FAILURE":
				if len(data) > 2 {
					errorMessage := strings.Join(data[2:], " ")
					return errors.New(errorMessage)
				}
			}
			break
		}
	}
	return errors.New("error registering user")
}

// ch: send request accepted or denied
func (cl *CentralClient) HandleSendRequest(handler func(p *packet.SendPacket) bool) {
	cl.sendRequestHandler = handler
}

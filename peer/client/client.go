package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"share/common/packet"
	"share/peer/protocol"
	"strings"
)

type CentralClient struct {
	Conn net.Conn
	// @returns if this handler has reached its end-case (removes from arr if true)
	dataHandlers       []func([]byte) bool
	sendRequestHandler func(*packet.SendPacket) bool
	Started            bool
	RegisteredUsername string
}

func NewClient() CentralClient {
	conn, err := net.Dial("tcp", "localhost:8080")
	catch(err)
	centralClient := CentralClient{
		Conn:               conn,
		dataHandlers:       []func([]byte) bool{},
		sendRequestHandler: nil,
		Started:            false,
		RegisteredUsername: "",
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

	_, err := fmt.Fprintln(cl.Conn, p.String())
	catch(err)
}

func (cl *CentralClient) HandleData(callback func([]byte) bool) {
	cl.dataHandlers = append(cl.dataHandlers, callback)
}

func (cl *CentralClient) RemoveDataHandler(index int) {
	last := len(cl.dataHandlers) - 1
	if last < 0 {
		return
	}
	cl.dataHandlers[last], cl.dataHandlers[index] = cl.dataHandlers[index], cl.dataHandlers[last]
}

func (cl *CentralClient) HandleError(errorType string, callback func(err string)) {
	handler := func(buf []byte) bool {
		data := string(buf)
		packetType := packet.GetPacketType(data)
		if packetType == packet.ERROR {
			p := packet.ToErrorPacket(data)
			if p.ErrorType == errorType {
				callback(p.ErrorMessage)
				return true
			}
		}
		return false
	}
	cl.HandleData(handler)
}

func (cl *CentralClient) handleData(buf []byte) {
	data := string(buf)
	fmt.Printf("I hath received: %s\n", data)
	packetType := packet.GetPacketType(data)
	if packetType == packet.SEND {
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

	for index, callback := range cl.dataHandlers {
		remove := callback(buf)
		if remove {
			cl.RemoveDataHandler(index)
		}
	}
}

func (cl *CentralClient) sendFileRequest(req *packet.SendPacket) (*packet.AcceptPacket, error) {
	cl.WritePacket(req)

	type Result struct {
		Packet *packet.AcceptPacket
		Error  error
	}
	ch := make(chan Result)

	callback := func(buf []byte) bool {
		data := string(buf)
		packetType := packet.GetPacketType(data)
		switch packetType {
		case packet.ACCEPT:
			acceptPacket := packet.ToAcceptPacket(data)
			if acceptPacket.Filename == req.Filename && acceptPacket.Size == req.Size {
				ch <- Result{
					Packet: acceptPacket,
					Error:  nil,
				}
				return true
			}
		case packet.REJECT:
			rejectPacket := packet.ToRejectPacket(data)
			if rejectPacket.Filename == req.Filename {
				ch <- Result{nil, nil}
				return true
			}
		case packet.ERROR:
			errorPacket := packet.ToErrorPacket(data)
			if errorPacket.ErrorType == packet.SEND {
				ch <- Result{Packet: nil, Error: fmt.Errorf(errorPacket.ErrorMessage)}
				return true
			}
		}
		return false
	}
	cl.HandleData(callback)
	result := <-ch
	return result.Packet, result.Error
}

// TODO: continue from here

// @returns (bytes sent, error)
func (cl *CentralClient) SendFile(filepath string, target string) (int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return -1, err
	}
	info, err := file.Stat()
	if err != nil {
		return -1, err
	}
	filesize := int(info.Size())

	app := protocol.NewShareHandler()
	req := packet.NewSendPacket(filepath, filesize, target, app.PeerHandler.GetPeerAddress())
	accept, err := cl.sendFileRequest(&req)
	if err != nil {
		return -1, err
	}
	if accept != nil {
		if err := app.Send(accept); err != nil {
			return filesize, err
		}
	}

	return filesize, nil
}

func (cl *CentralClient) RegisterUsername(username string) error {
	registerPacket := packet.NewRegisterPacket(username)
	cl.WritePacket(&registerPacket)

	ch := make(chan error)
	callback := func(buf []byte) bool {
		data := strings.Split(string(buf), " ")
		if len(data) < 2 {
			return false
		}
		if data[0] == registerPacket.PacketType {
			switch data[1] {
			case "USER_REGISTER_SUCCESS":
				cl.RegisteredUsername = username
				ch <- nil
			case "USER_REGISTER_FAILURE":
				if len(data) > 2 {
					errorMessage := strings.Join(data[2:], " ")
					ch <- fmt.Errorf(errorMessage)
				} else {
					ch <- fmt.Errorf("error registering user")
				}
			}
			return true
		}
		return false
	}
	cl.HandleData(callback)
	return <-ch
}

// ch: send request accepted or denied
func (cl *CentralClient) HandleSendRequest(handler func(p *packet.SendPacket) bool) {
	cl.sendRequestHandler = handler
}

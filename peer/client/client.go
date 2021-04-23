package client

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"share/common/packet"
	"share/peer/application"
	"strings"
)

type CentralClient struct {
	Conn               net.Conn
	DataChannels       []chan []byte
	sendRequestHandler func(p *packet.SendPacket) bool
	Started            bool
	App                *application.Application
}

func NewClient(app *application.Application) CentralClient {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	centralClient := CentralClient{
		Conn:         conn,
		DataChannels: []chan []byte{},
		Started:      false,
		App:          app,
	}
	return centralClient
}

func (cl *CentralClient) Start() {
	cl.Started = true
	defer func() {
		if err := cl.Conn.Close(); err != nil {
			log.Fatalf("Error: %s", err)
		}
	}()

	scanner := bufio.NewScanner(cl.Conn)
	for scanner.Scan() {
		go cl.handleData(scanner.Bytes())
	}
}

func (cl *CentralClient) WritePacket(p packet.Packet) {
	if !cl.Started {
		log.Fatal("Error: trying to write data before connected.")
		return
	}

	data := p.String()
	fmt.Printf("Writing %s\n", data)
	if _, err := fmt.Fprintln(cl.Conn, data); err != nil {
		log.Fatalf("Error: %s", err)
	}
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
	fmt.Println(data)
	packetType := packet.GetPacketType(data)
	if packetType == "FILE_SEND_REQUEST" {
		p := packet.ToSendPacket(data)
		res := cl.sendRequestHandler(p)
		if res {
			accepted := packet.NewAcceptPacket(p.Filename, p.Size, cl.App.GetPeerAddress())
			cl.WritePacket(&accepted)
		} else {
			rejected := packet.NewRejectPacket(p.Filename)
			cl.WritePacket(&rejected)
		}
	}

	for _, c := range cl.DataChannels {
		c <- buf
	}
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

package client

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"share/common/packet"
	"strings"
)

type CentralClient struct {
	Conn                net.Conn
	DataChannels        []chan []byte
	sendRequestHandler func(p *packet.SendPacket) bool
	PeerAddr string
}

func NewClient() CentralClient {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	return CentralClient{
		Conn:         conn,
		DataChannels: []chan []byte{},
	}
}

func (a *CentralClient) Start() {
	defer func() {
		if err := a.Conn.Close(); err != nil {
			log.Fatalf("Error: %s", err)
		}
	}()

	scanner := bufio.NewScanner(a.Conn)
	for scanner.Scan() {
		go a.handleData(scanner.Bytes())
	}
}

func (a *CentralClient) WriteData(data string) {
	if _, err := fmt.Fprintln(a.Conn, data); err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func (a *CentralClient) CreateDataChannel() chan []byte {
	ch := make(chan []byte)
	a.DataChannels = append(a.DataChannels, ch)
	return ch
}

func (a *CentralClient) RemoveDataChannel(c chan []byte) {
	for i, ch := range a.DataChannels {
		if ch == c {
			last := len(a.DataChannels) - 1
			a.DataChannels[i], a.DataChannels[last] = a.DataChannels[last], a.DataChannels[i]
			a.DataChannels = a.DataChannels[:last]
		}
	}
}

func (a *CentralClient) handleData(buf []byte) {
	data := string(buf)
	packetType := packet.GetPacketType(data)
	if packetType == "FILE_SEND_REQUEST" {
		p := packet.ToSendPacket(data)
		res := a.sendRequestHandler(p)
		if res {
			accepted := packet.NewAcceptPacket(p.Filename, p.Size, a.PeerAddr)
			a.WriteData(accepted.Serialize())
		} else {
			//rejected :=
		}
		return
	}

	for _, c := range a.DataChannels {
		c <- buf
	}
}

func (a *CentralClient) RegisterUsername(username string) error {
	registerPacket := packet.NewRegisterPacket(username)
	c := a.CreateDataChannel()
	a.WriteData(registerPacket.Serialize())
	for {
		res := <-c
		data := strings.Split(string(res), " ")
		if len(data) < 2 {
			continue
		}
		if data[0] == registerPacket.PacketType {
			a.RemoveDataChannel(c)
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
func (a *CentralClient) HandleSendRequest(handler func(p *packet.SendPacket) bool) {
	a.sendRequestHandler = handler
}

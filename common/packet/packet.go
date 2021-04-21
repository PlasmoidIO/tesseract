package packet

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	SEND      = "FILE_SEND_REQUEST"
	ACCEPT    = "SEND_REQUEST_ACCEPTED"
	REJECT    = "SEND_REQUEST_REJECTED"
	REGISTER  = "REGISTER_USERNAME"
	SEPARATOR = " "
)

type Packet interface {
	Serialize() string
}

type SendPacket struct {
	PacketType string
	Filename   string
	Size       int
	User       string
}

type AcceptPacket struct {
	PacketType string
	Filename   string
	Size       int
	PeerAddr   string
}

type RejectPacket struct {
	PacketType string
	Filename   string
}

type RegisterPacket struct {
	PacketType string
	Username   string
}

func ToRegisterPacket(data string) *RegisterPacket {
	arr := strings.Split(data, SEPARATOR)
	if len(arr) < 2 {
		return nil
	}
	return &RegisterPacket{
		PacketType: arr[0],
		Username:   arr[1],
	}
}

func ToSendPacket(data string) *SendPacket {
	arr := strings.Split(data, SEPARATOR)
	if len(arr) < 4 {
		return nil
	}
	n, err := strconv.Atoi(arr[2])
	if err != nil {
		return nil
	}
	return &SendPacket{
		PacketType: arr[0],
		Filename:   arr[1],
		Size:       n,
		User:       arr[3],
	}
}

func GetPacketType(packet string) string {
	arr := strings.Split(packet, SEPARATOR)
	if len(arr) > 0 {
		return arr[0]
	}
	return ""
}

func (p *RegisterPacket) Serialize() string {
	return fmt.Sprintf("%s%s%s", p.PacketType, SEPARATOR, p.Username)
}

func (p *AcceptPacket) Serialize() string {
	return fmt.Sprintf("%s%s%s%s%d%s%s", p.PacketType, SEPARATOR, p.Filename, SEPARATOR, p.Size, SEPARATOR, p.PeerAddr)
}

func (p *SendPacket) Serialize() string {
	return fmt.Sprintf("%s%s%s%s%d%s%s", p.PacketType, SEPARATOR, p.Filename, SEPARATOR, p.Size, SEPARATOR, p.User)
}

func (p *RejectPacket) Serialize() string {
	return fmt.Sprintf("%s%s%s", p.PacketType, SEPARATOR, p.Filename)
}

func NewRegisterPacket(username string) RegisterPacket {
	return RegisterPacket{
		PacketType: REGISTER,
		Username:   username,
	}
}

func NewSendPacket(filename string, size int, user string) SendPacket {
	return SendPacket{
		PacketType: SEND,
		Filename:   filename,
		Size:       size,
		User:       user,
	}
}

func NewAcceptPacket(filename string, size int, peeraddr string) AcceptPacket {
	return AcceptPacket{
		PacketType: ACCEPT,
		Filename:   filename,
		Size:       size,
		PeerAddr:   peeraddr,
	}
}

func NewRejectPacket(filename string) RejectPacket {
	return RejectPacket{
		PacketType: REJECT,
		Filename:   filename,
	}
}

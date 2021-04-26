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
	ERROR     = "ERROR"
	SEPARATOR = " "
)

type Packet interface {
	String() string
}

type SendPacket struct {
	PacketType string
	Filename   string
	Size       int
	Username   string
	SenderAddr string
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

type ErrorPacket struct {
	PacketType   string
	ErrorType    string
	ErrorMessage string
}

func ToErrorPacket(data string) *ErrorPacket {
	arr := strings.Split(data, SEPARATOR)
	if len(arr) < 3 {
		return nil
	}
	return &ErrorPacket{
		PacketType:   arr[0],
		ErrorType:    arr[1],
		ErrorMessage: strings.Join(arr[2:], " "),
	}
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
	if len(arr) < 5 {
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
		Username:   arr[3],
		SenderAddr: arr[4],
	}
}

func ToAcceptPacket(data string) *AcceptPacket {
	arr := strings.Split(data, SEPARATOR)
	if len(arr) < 4 {
		return nil
	}
	n, err := strconv.Atoi(arr[2])
	if err != nil {
		return nil
	}
	return &AcceptPacket{
		PacketType: arr[0],
		Filename:   arr[1],
		Size:       n,
		PeerAddr:   arr[3],
	}
}

func ToRejectPacket(data string) *RejectPacket {
	arr := strings.Split(data, SEPARATOR)
	if len(arr) < 2 {
		return nil
	}
	return &RejectPacket{
		PacketType: arr[0],
		Filename:   arr[1],
	}
}

func GetPacketType(packet string) string {
	arr := strings.Split(packet, SEPARATOR)
	if len(arr) > 0 {
		return arr[0]
	}
	return ""
}

func (p *RegisterPacket) String() string {
	return fmt.Sprintf("%s%s%s", p.PacketType, SEPARATOR, p.Username)
}

func (p *AcceptPacket) String() string {
	return fmt.Sprintf("%s%s%s%s%d%s%s", p.PacketType, SEPARATOR, p.Filename, SEPARATOR, p.Size, SEPARATOR, p.PeerAddr)
}

func (p *SendPacket) String() string {
	return fmt.Sprintf("%s%s%s%s%d%s%s%s%s", p.PacketType, SEPARATOR, p.Filename, SEPARATOR, p.Size, SEPARATOR, p.Username, SEPARATOR, p.SenderAddr)
}

func (p *RejectPacket) String() string {
	return fmt.Sprintf("%s%s%s", p.PacketType, SEPARATOR, p.Filename)
}

func (p *ErrorPacket) String() string {
	return fmt.Sprintf("%s%s%s%s%s", p.PacketType, SEPARATOR, p.ErrorType, SEPARATOR, p.ErrorMessage)
}

func NewRegisterPacket(username string) RegisterPacket {
	return RegisterPacket{
		PacketType: REGISTER,
		Username:   username,
	}
}

func NewSendPacket(filename string, size int, target string, senderAddr string) SendPacket {
	return SendPacket{
		PacketType: SEND,
		Filename:   filename,
		Size:       size,
		Username:   target,
		SenderAddr: senderAddr,
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

func NewErrorPacket(errorType string, errorMessage string) ErrorPacket {
	return ErrorPacket{
		PacketType:   ERROR,
		ErrorType:    errorType,
		ErrorMessage: errorMessage,
	}
}

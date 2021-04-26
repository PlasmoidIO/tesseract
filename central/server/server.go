package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"share/central/authentication"
	"share/common/packet"
)

type Server struct {
	channels map[chan string]net.Conn
	// username: connection
	authHandler authentication.AuthHandler
}

/**
  - Username registers
  - Username send file request [target(username), filename, filesize]
*/
func catch(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func write(w io.Writer, a interface{}) {
	_, err := fmt.Fprintln(w, a)
	catch(err)
}

func (s *Server) handleRegisterPacket(conn net.Conn, p *packet.RegisterPacket) {
	var response string
	if s.authHandler.LoginUser(conn, p.Username) {
		response = fmt.Sprintf("%s USER_REGISTER_SUCCESS", p.PacketType)
	} else {
		response = fmt.Sprintf("%s %s Name already taken.", packet.ERROR, p.PacketType)
	}
	write(conn, response)
}

func (s *Server) handleSendPacket(conn net.Conn, p *packet.SendPacket) {
	user, has := s.authHandler.Connected[p.Username]
	if has {
		sender, has := s.getUserFromConn(conn)
		if has {
			p.Username = sender
			write(user, p.String())

			ch := make(chan string)
			s.channels[ch] = user
			response := <-ch
			delete(s.channels, ch)
			write(conn, response)
		} else {
			write(conn, fmt.Sprintf("%s %s NOT_AUTHORIZED", packet.ERROR, p.PacketType))
		}
		return
	}
	write(conn, fmt.Sprintf("%s %s USER_NOT_FOUND", packet.ERROR, p.PacketType))
}

func (s *Server) handleResponsePacket(conn net.Conn, data string) {
	for ch, c := range s.channels {
		if c == conn {
			ch <- data
		}
	}
}

func (s *Server) handleData(conn net.Conn, buf []byte) {
	data := string(buf)
	switch packet.GetPacketType(data) {
	case packet.REGISTER:
		p := packet.ToRegisterPacket(data)
		if p == nil {
			break
		}
		s.handleRegisterPacket(conn, p)
	case packet.SEND:
		p := packet.ToSendPacket(data)
		if p == nil {
			break
		}
		s.handleSendPacket(conn, p)
	case packet.ACCEPT, packet.REJECT:
		s.handleResponsePacket(conn, data)
	}
}

func (s *Server) getUserFromConn(conn net.Conn) (string, bool) {
	for username, c := range s.authHandler.Connected {
		if c == conn {
			return username, true
		}
	}
	return "", false
}

// returns success

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		for k, v := range s.authHandler.Connected {
			if conn == v {
				delete(s.authHandler.Connected, k)
				break
			}
		}
		if err := conn.Close(); err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		buf := scanner.Bytes()
		go s.handleData(conn, buf)
	}
}

func CreateServer() Server {
	return Server{
		channels:    make(map[chan string]net.Conn),
		authHandler: authentication.NewHandler(),
	}
}

func (s *Server) Listen(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	fmt.Printf("Listening on port %d\n", port)
	catch(err)
	for {
		conn, err := listener.Accept()
		catch(err)
		go s.handleConnection(conn)
	}
}

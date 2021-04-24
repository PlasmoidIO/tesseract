package application

import (
	"github.com/libp2p/go-libp2p-core/network"
	"io/ioutil"
	"log"
	"share/common/packet"
)

type ShareHandler struct {
	PeerHandler PeerHandler
}

func NewShareHandler() ShareHandler {
	peerHandler := NewPeerHandler()
	return ShareHandler{peerHandler}
}

func (s *ShareHandler) Send(req *packet.AcceptPacket) {
	stream := s.PeerHandler.OpenConnection(req.PeerAddr)
	defer func() {
		catch(stream.Close())
	}()
	buf, err := ioutil.ReadFile(req.Filename)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = stream.Write(buf); err != nil {
		log.Fatalf("Error writing to stream: %s", err)
	}
}

// TODO: stream authentication
func (s *ShareHandler) Receive(req *packet.SendPacket) {
	ch := make(chan bool)
	callback := func(stream network.Stream) {
		defer func() {
			catch(stream.Close())
		}()
		buf, err := ioutil.ReadAll(stream)
		if err != nil {
			log.Fatalf("Error reading from stream: %s", err)
		}
		if err := ioutil.WriteFile(req.Filename, buf, 0); err != nil {
			log.Fatalf("Error writing to file: %s", err)
		}
		ch <- true
	}
	s.PeerHandler.HandleIncoming(callback)
	<-ch
}

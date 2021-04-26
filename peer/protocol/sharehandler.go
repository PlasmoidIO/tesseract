package protocol

import (
	"fmt"
	"io/ioutil"
	"share/common/packet"

	"github.com/libp2p/go-libp2p-core/network"
)

type ShareHandler struct {
	PeerHandler PeerHandler
}

func NewShareHandler() ShareHandler {
	peerHandler := NewPeerHandler()
	return ShareHandler{peerHandler}
}

func (s *ShareHandler) Send(req *packet.AcceptPacket) error {
	stream := s.PeerHandler.OpenConnection(req.PeerAddr)
	defer func() {
		catch(stream.Close())
	}()
	buf, err := ioutil.ReadFile(req.Filename)
	catch(err)
	if _, err = stream.Write(buf); err != nil {
		return fmt.Errorf("error writing to stream: %s", err)
	}
	return nil
}

// TODO: stream authentication
func (s *ShareHandler) Receive(req *packet.SendPacket) error {
	ch := make(chan error)

	callback := func(stream network.Stream) {
		defer func() {
			catch(stream.Close())
		}()
		buf, err := ioutil.ReadAll(stream)
		if err != nil {
			ch <- fmt.Errorf("error reading from stream: %s", err)
			return
		}
		if err := ioutil.WriteFile(req.Filename, buf, 0); err != nil {
			ch <- fmt.Errorf("error writing to file: %s", err)
			return
		}
		ch <- nil
	}

	s.PeerHandler.HandleIncoming(callback)
	return <-ch
}

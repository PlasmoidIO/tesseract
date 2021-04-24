package application

import (
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p-core/network"
	"io/ioutil"
	"share/common/packet"
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
		return errors.New(fmt.Sprintf("Error writing to stream: %s", err))
	}
	return nil
}

// TODO: stream authentication
func (s *ShareHandler) Receive(req *packet.SendPacket) error {
	ch := make(chan string)
	callback := func(stream network.Stream) {
		defer func() {
			catch(stream.Close())
		}()
		buf, err := ioutil.ReadAll(stream)
		if err != nil {
			ch <- fmt.Sprintf("Error reading from stream: %s", err)
			return
		}
		if err := ioutil.WriteFile(req.Filename, buf, 0); err != nil {
			ch <- fmt.Sprintf("Error writing to file: %s", err)
			return
		}
		ch <- ""
	}
	s.PeerHandler.HandleIncoming(callback)
	res := <-ch
	if res != "" {
		return errors.New(res)
	}
	return nil
}

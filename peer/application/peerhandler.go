package application

import (
	"context"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"log"
)

const PROTOCOL = "share"

type PeerHandler struct {
	host     host.Host
	peerAddr string
}

func NewPeerHandler() PeerHandler {
	ctx := context.Background()
	node, err := libp2p.New(ctx)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	peerInfo := peer.AddrInfo{ID: node.ID(), Addrs: node.Addrs()}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	return PeerHandler{
		host:     node,
		peerAddr: addrs[0].String(),
	}
}

func (p *PeerHandler) Close() {
	if err := p.host.Close(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func (p *PeerHandler) GetPeerAddress() string {
	return p.peerAddr
}

func (p *PeerHandler) HandleIncoming(callback func(stream network.Stream)) {
	p.host.SetStreamHandler(PROTOCOL, callback)
}

func (p *PeerHandler) OpenConnection(address string) network.Stream {
	ma, err := multiaddr.NewMultiaddr(address)
	catch(err)
	addrInfo, err := peer.AddrInfoFromP2pAddr(ma)
	catch(err)
	catch(p.host.Connect(context.Background(), *addrInfo))
	stream, err := p.host.NewStream(context.Background(), addrInfo.ID, PROTOCOL)
	catch(err)
	return stream
}

package protocol

import (
	"context"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

const PROTOCOL = "share"

type PeerHandler struct {
	host     host.Host
	peerAddr string
}

func NewPeerHandler() PeerHandler {
	ctx := context.Background()
	node, err := libp2p.New(ctx)
	catch(err)
	peerInfo := peer.AddrInfo{ID: node.ID(), Addrs: node.Addrs()}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	catch(err)
	return PeerHandler{
		host:     node,
		peerAddr: addrs[0].String(),
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

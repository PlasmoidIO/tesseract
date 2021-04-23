package application

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"log"
	"share/common/packet"
	"share/peer/client"
)

const PROTOCOL = "share"

type Application struct {
	host     host.Host
	peerAddr string

	Client client.CentralClient
}

func (a *Application) GetPeerAddress() string {
	return a.peerAddr
}

func NewApplication() Application {
	ctx := context.Background()
	node, err := libp2p.New(ctx)
	catch(err)
	app := Application{host: node}
	cl := client.NewClient(&app)
	app.Client = cl
	return app
}

func catch(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func (a *Application) handleSendRequest(req *packet.SendPacket) bool {
	fmt.Printf("Send request: %s\n", req)
	return true
}

func (a *Application) handleConnection(stream network.Stream) {

}

func (a *Application) Start() {
	defer func() {
		err := a.host.Close()
		catch(err)
	}()

	a.Client.HandleSendRequest(a.handleSendRequest)
	a.host.SetStreamHandler(PROTOCOL, func(s network.Stream) {
		go a.handleConnection(s)
	})

	peerInfo := peer.AddrInfo{
		ID:    a.host.ID(),
		Addrs: a.host.Addrs(),
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	catch(err)
	a.peerAddr = addrs[0].String()

	a.Client.Start()
}

package client

import (
	"context"
	"fmt"
	"net/http"

	"libp2p-dht-discover/log"

	dht "github.com/libp2p/go-libp2p-kad-dht"

	// log "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p/core/event"
	"github.com/libp2p/go-libp2p/core/host"
)

type P2PClient struct {
	node       host.Host
	dht        *dht.IpfsDHT
	httpServer *http.Server
	ctx        context.Context
}

func NewClient(ctx context.Context, node host.Host) *P2PClient {
	return &P2PClient{
		node: node,
		ctx:  ctx,
	}
}

// logs the addresses of the host every time they change
func (c *P2PClient) ListenForAddrChanges() {
	// Subscribe to address change events
	sub, err := c.node.EventBus().Subscribe(new(event.EvtLocalAddressesUpdated))
	if err != nil {
		log.Fatalf("Failed to subscribe to address change events: %s", err)
	}
	defer sub.Close()

	for {
		select {
		case evt := <-sub.Out():
			if e, ok := evt.(event.EvtLocalAddressesUpdated); ok {
				log.Info("Addresses have been updated!")
				for _, addr := range e.Current {
					fullAddr := fmt.Sprintf("%+v/p2p/%s", addr, c.node.ID().Pretty())
					log.Infof("Advertised Address:", fullAddr)
				}
			}
		case <-c.ctx.Done():
			return
		}
	}
}

package client

import (
	"libp2p-dht-discover/log"

	kadht "github.com/libp2p/go-libp2p-kad-dht"
)

// list of default bootstrap peers
// var DefaultBootstrapPeers = []string{
// 	"/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
// 	"/ip4/104.236.176.52/tcp/4001/p2p/QmSoLnSGccFuZQJzRadHn95W2CrSFmZuTdDWP8HXaHca9z",
// 	"/ip4/104.236.179.241/tcp/4001/p2p/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",
// 	"/ip4/162.243.248.213/tcp/4001/p2p/QmSoLueR4xBeUbY9WZ9xGUUxunbKWcrNFTDAadQJmocnWm",
// 	"/ip4/128.199.219.111/tcp/4001/p2p/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",
// 	"/ip4/104.236.76.40/tcp/4001/p2p/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64",
// 	"/ip4/178.62.158.247/tcp/4001/p2p/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd",
// 	"/ip4/178.62.61.185/tcp/4001/p2p/QmSoLMeWqB7YGVLJN3pNLQpmmEk35v6wYtsMGLzSr5QBU3",
// 	"/ip4/104.236.151.122/tcp/4001/p2p/QmSoLju6m7xTh3DuokvT3886QRYqxAzb1kShaanJgW36yx",
// }

// initDHT initializes the DHT and bootstrap it with the default bootstrap peers
func (c *P2PClient) InitDHT() error {
	// create a new DHT
	dht, err := kadht.New(c.ctx, c.node)
	if err != nil {
		log.Fatalf("error creating DHT: %s", err)
	}
	// attach the DHT to the libp2p host
	c.dht = dht

	defaultBootstrapAddInfo := kadht.GetDefaultBootstrapPeerAddrInfos()
	for _, peerInfo := range defaultBootstrapAddInfo {
		log.Debugf("adding bootstraping node: %s", peerInfo.ID.Pretty())
		err := c.node.Connect(c.ctx, peerInfo)
		if err != nil {
			log.Errorf("error connecting to bootstrap node: %s", err)
			continue
		}
		log.Debugf("connected to bootstrap node: %s", peerInfo.ID.Pretty())

	}

	// Bootstrap the DHT
	return dht.Bootstrap(c.ctx)

}

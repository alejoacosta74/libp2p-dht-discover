package node

import (
	"fmt"

	"libp2p-dht-discover/log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"

	multiaddr "github.com/multiformats/go-multiaddr"
)

// name of the file that contains the Docker's host IP
const hostIPFile = "hostip"

func NewNode(ipaddr string, port string) (host.Host, error) {

	privateKey, err := getPrivKey()
	if err != nil {
		log.Fatalf("error getting private key: %s", err)
	}

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", ipaddr, port))

	// check if there is a host IP we can use
	hostIP, err := readHostIPFromFile(hostIPFile)
	if err != nil || hostIP == "" {
		log.Errorf("error reading host IP from file: %s. Skiping...", err)
		node, err := libp2p.New(
			libp2p.Identity(privateKey),
			libp2p.ListenAddrs(sourceMultiAddr),
			libp2p.EnableNATService(),
			libp2p.NATPortMap(),
		)
		if err != nil {
			log.Errorf("error creating node: %s", err)
			return nil, err
		}
		return node, nil
	}

	log.Infof("found host IP: %s", hostIP)
	addrsFactory := makeAddrsFactory(hostIP)
	node, err := libp2p.New(
		libp2p.Identity(privateKey),
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.AddrsFactory(addrsFactory),
		libp2p.EnableNATService(),
		libp2p.NATPortMap(),
	)
	if err != nil {
		log.Errorf("error creating node: %s", err)
		return nil, err
	}
	return node, nil
}

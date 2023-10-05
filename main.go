package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	// "github.com/libp2p/go-libp2p-kad-dht"

	// log "github.com/ipfs/go-log/v2"
	"libp2p-dht-discover/client"
	"libp2p-dht-discover/log"
	"libp2p-dht-discover/utils"

	"github.com/libp2p/go-libp2p"

	multiaddr "github.com/multiformats/go-multiaddr"
)

func main() {
	fmt.Println("starting...")
	// log.SetAllLoggers(log.LevelInfo)
	log.ConfigureLogger(log.WithLevel("debug"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	privateKey, err := utils.GetPrivKey()
	if err != nil {
		log.Fatalf("error getting private key: %s", err)
	}

	ipaddr := "0.0.0.0"
	port := "2001"
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", ipaddr, port))

	node, err := libp2p.New(
		libp2p.Identity(privateKey),
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.EnableNATService(),
		libp2p.NATPortMap(),
	)
	if err != nil {
		log.Fatalf("error creating node: %s", err)
	}

	// log the node's listening addresses
	for _, addr := range node.Addrs() {
		log.Infof("listening on: %s", addr)
	}

	client := client.NewClient(ctx, node)

	log.Infof("node created: %s", node.ID().Pretty())
	// start the http server
	go func() {
		if err := client.StartServer(); err != nil {
			log.Fatalf("error starting http server: %s", err)
			os.Exit(1)
		}
	}()
	// Start listening for address changes
	go client.ListenForAddrChanges()

	// initialize and bootstrap the DHT
	if err := client.InitDHT(); err != nil {
		log.Fatalf("error initializing DHT: %s", err)
		os.Exit(1)
	}

	// wait for a SIGINT or SIGTERM signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Warn("Received signal, shutting down...")
	if err := node.Close(); err != nil {
		panic(err)
	}
}

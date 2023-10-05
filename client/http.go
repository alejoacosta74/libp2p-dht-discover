package client

import (
	"context"
	"libp2p-dht-discover/log"
	"net/http"
	"strings"

	// log "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p/core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func (c *P2PClient) StartServer() error {
	// start http server
	mux := c.newServer()
	c.httpServer = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	return c.httpServer.ListenAndServe()
}

func (c *P2PClient) StopServer() {
	err := c.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("error stopping http server: %s", err)
	}
}

// handler for the /dhtpeers endpoint
func (c *P2PClient) dhtPeersHandler(w http.ResponseWriter, r *http.Request) {
	// get the list of peers from the dht
	peers := c.dht.RoutingTable().ListPeers()
	log.Debugf("peers on dht: %d", len(peers))
	// write the list of peers to the response
	for _, peer := range peers {
		w.Write([]byte(peer.Pretty() + "\n"))
	}
}

// handler for the /nodepeers endpoint
func (c *P2PClient) nodePeersHandler(w http.ResponseWriter, r *http.Request) {
	// get the list of peers from the host
	peers := c.node.Peerstore().Peers()
	// write the list of peers to the response
	log.Debugf("peers on peerStore: %d", len(peers))
	for _, peer := range peers {
		w.Write([]byte(peer.Pretty() + "\n"))
	}
}

// handler for the /connect/{multiaddr} endpoint
func (c *P2PClient) connectHandler(w http.ResponseWriter, r *http.Request) {
	addr := strings.TrimPrefix(r.URL.Path, "/connect")
	maddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		log.Errorf("Failed to parse multiaddress: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	addrInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Errorf("Failed to parse multiaddress: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Discover the Mac node's addresses using DHT
	discoveredAddrs, err := c.dht.FindPeer(c.ctx, addrInfo.ID)
	if err != nil {
		log.Errorf("Failed to discover peer: %v, error: %s", addrInfo.ID, err)
		http.Error(w, "Failed to discover peer: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Discovered addresses:", discoveredAddrs)

	ctx := context.Background()
	if err := c.node.Connect(ctx, *addrInfo); err != nil {
		log.Errorf("Failed to connect to peer: %v, %s", addrInfo.ID.Pretty(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Connected successfully to peer " + addrInfo.ID.Pretty()))
}

// NewServer creates a new http server with the available handlers
func (c *P2PClient) newServer() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/dhtpeers", c.dhtPeersHandler)
	mux.HandleFunc("/nodepeers", c.nodePeersHandler)
	mux.HandleFunc("/connect/", c.connectHandler)
	return mux
}

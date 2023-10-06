package node

import (
	"libp2p-dht-discover/log"
	"os"
	"strings"

	multiaddr "github.com/multiformats/go-multiaddr"
)

// returns a function that replaces the internal IP (i.e. Docker IP) with the host IP
func makeAddrsFactory(hostIP string) func([]multiaddr.Multiaddr) []multiaddr.Multiaddr {
	return func(addrs []multiaddr.Multiaddr) []multiaddr.Multiaddr {
		var newAddrs []multiaddr.Multiaddr
		var dockerIP string

		// Find the first non-loopback IP, which should be the Docker container's IP
		for _, addr := range addrs {
			if !isLoopbackAddr(addr) {
				parts := strings.Split(addr.String(), "/")
				if len(parts) > 1 && parts[1] == "ip4" {
					log.Debugf("found docker IP: %s", parts[2])
					dockerIP = parts[2]
					break
				}
			}
		}

		// Replace the Docker IP with the host IP
		for _, addr := range addrs {
			log.Debugf("replacing docker IP with address: %s", addr.String())
			if dockerIP != "" {
				newAddr, err := multiaddr.NewMultiaddr(strings.Replace(addr.String(), dockerIP, hostIP, 1))
				if err == nil {
					newAddrs = append(newAddrs, newAddr)
				} else {
					log.Errorf("error creating new multiaddr from addr %sm errir: %s", addr.String(), err)
				}
			} else {
				newAddrs = append(newAddrs, addr)
			}
		}
		return newAddrs
	}
}

func isLoopbackAddr(addr multiaddr.Multiaddr) bool {
	return strings.Contains(addr.String(), "127.0.0.1") || strings.Contains(addr.String(), "localhost")
}

// reads the host IP from a file
func readHostIPFromFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

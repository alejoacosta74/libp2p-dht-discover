package node

import (
	"crypto/rand"
	"os"

	"libp2p-dht-discover/log"

	"github.com/libp2p/go-libp2p/core/crypto"
)

// Returns the private key stored in the file to be used to generate
// the node's identity (CID)
func getPrivKey() (crypto.PrivKey, error) {
	// check if there is a file with the private key. It not, create one
	_, err := os.Stat("private.key")
	if os.IsNotExist(err) {
		// generate private key
		privateKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
		if err != nil {
			log.Errorf("error generating private key: %s", err)
			return nil, err

		}

		// save private key
		privateKeyBytes, err := crypto.MarshalPrivateKey(privateKey)
		if err != nil {
			log.Errorf("error marshalling private key: %s", err)
			return nil, err
		}
		err = os.WriteFile("private.key", privateKeyBytes, 0600)
		if err != nil {
			log.Errorf("error saving private key: %s", err)
			return nil, err
		}
	}

	// load private key
	privateKeyBytes, err := os.ReadFile("private.key")
	if err != nil {
		log.Errorf("error reading private key: %s", err)
		return nil, err
	}
	privateKey, err := crypto.UnmarshalPrivateKey(privateKeyBytes)
	if err != nil {
		log.Errorf("error unmarshalling private key: %s", err)
		return nil, err
	}
	return privateKey, nil

}

package internal

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

func GetKeyPair(key string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	// first check if key is hex
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, fmt.Errorf("error casting public key to ECDSA")
	}

	return privateKey, publicKeyECDSA, nil
}

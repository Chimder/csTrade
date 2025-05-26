package internal

import (
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
)

func Bit() {
	privKey, _ := btcec.NewPrivateKey()
	pubKey := privKey.PubKey().SerializeCompressed()
	pubKeyHash := btcutil.Hash160(pubKey)
	versionedPayload := append([]byte{0x00}, pubKeyHash...)

	firstHash := sha256.Sum256(versionedPayload)
	secondHash := sha256.Sum256(firstHash[:])
	checksum := secondHash[:4]

	fullPayload := append(versionedPayload, checksum...)
	address := base58.Encode(fullPayload)

	fmt.Printf("priv: %x\n", privKey.Serialize())
	fmt.Printf("pub key (hex): %x\n", pubKey)
	fmt.Printf("addr: %s\n", address)
}

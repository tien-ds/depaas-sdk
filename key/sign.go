package key

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
	"strings"
)

func NewSigner(base64PrivateKey string) auth.Signer {
	bytes, err := encoder.DecodeString(base64PrivateKey)
	if err != nil {
		panic(err)
	}
	signer := auth.NewSigner(auth.SignerTypeEd25519, bytes)
	return signer
}

func NewIdentity(mainNetPrivateKey string) *client.Identity {
	mainNetPrivKey, _ := crypto.HexToECDSA(strings.TrimPrefix(mainNetPrivateKey, "0x"))
	alice := &client.Identity{
		MainnetPrivKey: mainNetPrivKey,
	}
	return alice
}

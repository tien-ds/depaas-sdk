package key

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/loomnetwork/go-loom"
	"io"
	"signer/utils"
	"strings"
)

var encoder = base64.StdEncoding

func PrivateKeyToAddr(base64PrivateKey string) string {
	var priKey []byte
	if strings.HasPrefix(base64PrivateKey, "0x") {
		priKey, _ = hex.DecodeString(base64PrivateKey[2:])
	} else {
		priKey = []byte(base64PrivateKey)
	}
	bytes, err := encoder.DecodeString(string(priKey))
	if err != nil {
		panic(err)
	}
	var p = ed25519.PrivateKey(bytes)
	pub := p.Public().(ed25519.PublicKey)
	//fmt.Println(encoder.EncodeToString(pub[:]))
	addr := loom.LocalAddressFromPublicKey(pub)
	//fmt.Println(encoder.EncodeToString(addr))
	res := map[string]string{"addr": encoder.EncodeToString(addr), "pub": encoder.EncodeToString(pub[:])}
	byts, err := json.Marshal(res)
	utils.Require(err == nil, err)
	return string(byts)
}

func GenKey() *Key {
	return fromGenKey(nil)
}

func fromGenKey(rand io.Reader) *Key {
	pub, pri, err := ed25519.GenerateKey(rand)
	if err != nil {
		return nil
	}
	pubKeyB64 := encoder.EncodeToString(pub[:])
	privKeyB64 := encoder.EncodeToString(pri[:])
	addr := loom.LocalAddressFromPublicKey(pub[:])
	return &Key{
		Pub:  pubKeyB64,
		Pri:  privKeyB64,
		Addr: encoder.EncodeToString(addr),
	}
}

func FromSeed(seedHash string) *Key {
	buf, e := hex.DecodeString(seedHash)
	utils.Require(e == nil && len(buf) >= 32, "seed hash >=64")
	return fromGenKey(bytes.NewBuffer(buf[:32]))
}

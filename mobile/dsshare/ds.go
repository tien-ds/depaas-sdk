package dsshare

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"signer/key"
	"signer/rpc"
	"signer/token/ds"
)

var cDS *ds.DsContract

func SetDSContract(contract string) {
	cDS = ds.NewDsContract(rpc.GetDAppClient(), contract)
}

func DsStakeOf(addr string) string {
	return cDS.StackOf(addr).String()
}

func DsGetMinerAddr(peerId string) string {
	return base64.StdEncoding.EncodeToString(cDS.GetMinerAddr(peerId).Bytes())
}

func DsOwner() string {
	return base64.StdEncoding.EncodeToString(cDS.Owner().Bytes())
}

func DsPenalty(mainPrivateKey string, account string, amount string) string {
	a, b := big.NewInt(0).SetString(amount, 10)
	if !b {
		panic(fmt.Errorf("amount is error  %s", amount))
	}
	hash, err := cDS.Penalty(key.NewSigner(mainPrivateKey), account, a)
	if err != nil {
		panic(err)
	}
	return hash
}

func DsExchange(mainPrivateKey, sender, recipient, amount string) string {
	a, b := big.NewInt(0).SetString(amount, 10)
	if !b {
		panic(fmt.Errorf("amount is error  %s", amount))
	}
	exchange, err := cDS.Exchange(key.NewSigner(mainPrivateKey), sender, recipient, a)
	if err != nil {
		panic(err)
	}
	return exchange
}

func AddStake(mainPrivateKey string, sender string, amount string) string {
	am, b := big.NewInt(0).SetString(amount, 10)
	if !b {
		panic("amount is error")
	}
	stake, err := cDS.AddStake(key.NewSigner(mainPrivateKey), sender, am)
	if err != nil {
		panic(err)
	}
	return stake
}

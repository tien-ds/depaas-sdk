package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"signer/gateway"
	"signer/key"
	"signer/plugin"
	"signer/rpc"
	"signer/token/lm"
	"signer/utils"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/vm"
)

var oraclePools map[string]*gateway.Swap

func NewOracle(sMainNetCoin string, sLoomCoin string) *gateway.Swap {
	swap := gateway.NewSwap(rpc.GetMainNetClient(), rpc.GetDAppClient())
	swap.SetContracts(sMainNetCoin, sLoomCoin)
	return swap
}

func md5Hash(s string) string {
	bytes := md5.Sum([]byte(s))
	sML := fmt.Sprintf("%x", bytes)
	return sML
}

func hexToAddr(sHex string) string {
	bytes, err := hex.DecodeString(strings.TrimPrefix(sHex, "0x"))
	utils.Require(err == nil, err)
	return base64.StdEncoding.EncodeToString(bytes)
}

func AddrToHex(base64Addr string) string {
	bys, err := base64.StdEncoding.DecodeString(base64Addr)
	utils.Require(err == nil, err)
	return hex.EncodeToString(bys)
}

func genKey() string {
	bys, _ := json.Marshal(key.GenKey())
	return string(bys)
}

func genKeyFromSeed(seedHash string) string {
	bys, _ := json.Marshal(key.FromSeed(seedHash))
	return string(bys)
}

func getBalance(contractAddr string, hexAddr string) string {
	return rpc.GetBalance(contractAddr, hexAddr).String()
}

//{"contractAddr":"0xdb834d1f5baf312424fe3003524e2f5a52bf15b2","to":"0xdb834d1f5baf312424fe3003524e2f5a52bf15b2","base64PrivateKey":"","amount":10000}
func transfer(sJson string) string {
	st := struct {
		ContractAddr     string  `json:"contractAddr"`
		To               string  `json:"to"`
		Base64PrivateKey string  `json:"base64PrivateKey"`
		Amount           big.Int `json:"amount"`
	}(struct {
		ContractAddr     string
		To               string
		Base64PrivateKey string
		Amount           big.Int
	}{})
	err := json.Unmarshal([]byte(sJson), &st)
	utils.Require(err == nil, err)
	bys, err := json.Marshal(rpc.Transfer(struct {
		ContractAddr     string `json:"contractAddr"`
		To               string `json:"to"`
		Base64PrivateKey string `json:"base64PrivateKey"`
		Amount           string `json:"amount"`
	}(struct {
		ContractAddr     string
		To               string
		Base64PrivateKey string
		Amount           string
	}{ContractAddr: st.ContractAddr, To: st.To, Base64PrivateKey: st.Base64PrivateKey, Amount: st.Amount.String()})))
	utils.Require(err == nil, err)
	return string(bys)
}

func transfer1(sJson string) string {
	st := struct {
		ContractAddr     string `json:"contractAddr"`
		To               string `json:"to"`
		Base64PrivateKey string `json:"base64PrivateKey"`
		Amount           string `json:"amount"`
	}(struct {
		ContractAddr     string
		To               string
		Base64PrivateKey string
		Amount           string
	}{})
	err := json.Unmarshal([]byte(sJson), &st)
	utils.Require(err == nil, err)
	bys, err := json.Marshal(rpc.Transfer(st))
	utils.Require(err == nil, err)
	return string(bys)
}

func ethTransfer(url string, sjson string) string {
	p := struct {
		To      string  `json:"to"`
		Amount  big.Int `json:"amount"`
		ChainId big.Int `json:"chainId"`
		PriKey  string  `json:"priKey"`
	}{}
	err := json.Unmarshal([]byte(sjson), &p)
	utils.Require(err == nil, err)
	toAddr := common.HexToAddress(p.To)
	return rpc.EthTransfer(url, &toAddr, &p.Amount, &p.ChainId, p.PriKey)
}

//{"type":1,"hexAddr":"0xdb834d1f5baf312424fe3003524e2f5a52bf15b2"}
func getMapper(sJson string) string {
	st := struct {
		Type    int    `json:"type"`
		HexAddr string `json:"hexAddr"`
	}(struct {
		Type    int
		HexAddr string
	}{})
	err := json.Unmarshal([]byte(sJson), &st)
	utils.Require(err == nil, err)
	res := rpc.GetMapper(st)
	bys, err := json.Marshal(res)
	return string(bys)
}

func mapAccount(mainNetPrivateKey string, base64PrivateKey string) string {
	mapper := rpc.GetAccountMP()
	err := mapper.MapAccounts(mainNetPrivateKey, base64PrivateKey)
	return utils.ToJson(rpc.ToResult(err))
}

func tokenInfo(contractAddr string) string {
	dAppClient := rpc.GetDAppClient()
	la, _ := loom.LocalAddressFromHexString(utils.CosAddr(contractAddr))
	adc := lm.NewDAppChainContract(dAppClient, lm.ERC20ABI, loom.Address{
		ChainID: dAppClient.GetChainID(),
		Local:   la,
	})
	contr := lm.DAppChainEVMContract{ContractConnect: adc}
	return utils.ToJson(struct {
		Name     string `json:"name"`
		Decimals uint8  `json:"decimals"`
		Symbol   string `json:"symbol"`
	}{
		Name:     contr.Name(),
		Decimals: contr.Decimals(),
		Symbol:   contr.Symbol(),
	})

}

// eth => cos
func deposit(key string, mainNetPrivateKey, tokenAmount string) string {
	if swap, b := oraclePools[key]; b {
		amount := big.NewInt(0)
		amount.SetString(tokenAmount, 10)
		return utils.ToJson(rpc.ToResult(swap.Deposit(mainNetPrivateKey, amount)))
	}
	panic("deposit ptr is nil")

}

//cos => eth
func withdrawal(key string, base64PrivateKey string, mainNetPrivateKey string, tokenAmount string) string {
	if swap, b := oraclePools[key]; b {
		amount := big.NewInt(0)
		amount.SetString(tokenAmount, 10)
		return utils.ToJson(rpc.ToResult(swap.Withdrawal(base64PrivateKey, mainNetPrivateKey, amount)))
	}
	panic("withdrawal ptr is nil")

}

func newContracts(sMainNetCoin string, sLoomCoin string) string {
	hash := md5Hash(sMainNetCoin + sLoomCoin)
	if oraclePools == nil {
		oraclePools = make(map[string]*gateway.Swap)
	}
	if _, b := oraclePools[hash]; b {
		return hash
	} else {
		ptr := NewOracle(sMainNetCoin, sLoomCoin)
		oraclePools[hash] = ptr
		return hash
	}
}

func sign(base64key string, data []byte) string {
	keyBytes, err := base64.StdEncoding.DecodeString(base64key)
	utils.Require(err == nil, err)
	loomSigner := auth.NewEd25519Signer(keyBytes)
	//fmt.Println(base64.StdEncoding.EncodeToString(loomSigner.PublicKey()))
	return hex.EncodeToString(loomSigner.Sign(data))
}

func CommitCallTx(contract string, privateKey string, input []byte) string {
	lContract, err := loom.LocalAddressFromHexString(utils.CosAddr(contract))
	utils.Require(err == nil, err)
	dAppClient := rpc.GetDAppClient()
	newSigner := key.NewSigner(privateKey)
	bys, err := dAppClient.CommitCallTx(loom.Address{
		ChainID: dAppClient.GetChainID(),
		Local:   loom.LocalAddressFromPublicKey(newSigner.PublicKey()),
	}, loom.Address{
		ChainID: dAppClient.GetChainID(),
		Local:   lContract,
	}, newSigner, vm.VMType_EVM, input)
	utils.Require(err == nil, err)
	return hex.EncodeToString(bys)
}

func GateWayInit() {
	rpc.GateWayInit()
}

func DSInit() {
	rpc.InitClient()
}

func DSCoinBalance(account string) string {
	ethCoin, err := rpc.GetDAppClient().Resolve("ethcoin")
	if err != nil {
		panic(err)
	}
	coin := plugin.NewCoinPlugin(rpc.GetDAppClient(), ethCoin.Local)
	balance, err := coin.BalanceOf(account)
	if err != nil {
		panic(err)
	}
	return balance.String()
}

func DSCoinTransfer(base64PrivateKey string, to string, amount string) string {
	ethCoin, err := rpc.GetDAppClient().Resolve("ethcoin")
	if err != nil {
		panic(err)
	}
	coin := plugin.NewCoinPlugin(rpc.GetDAppClient(), ethCoin.Local)
	sKey := key.NewSigner(base64PrivateKey)
	bAmount, b := big.NewInt(0).SetString(amount, 10)
	if !b {
		panic(errors.New("amount error"))
	}
	hash, err := coin.Transfer(
		sKey,
		to,
		bAmount)
	if err != nil {
		panic(err)
	}
	return hash
}

func DSCoinMint(base64PrivateKey string, to string, amount string) string {
	ethCoin, err := rpc.GetDAppClient().Resolve("ethcoin")
	if err != nil {
		panic(err)
	}
	coin := plugin.NewCoinPlugin(rpc.GetDAppClient(), ethCoin.Local)
	sKey := key.NewSigner(base64PrivateKey)
	bAmount, b := big.NewInt(0).SetString(amount, 10)
	if !b {
		panic(errors.New("amount error"))
	}
	hash, err := coin.Mint(
		sKey,
		to,
		bAmount)
	if err != nil {
		panic(err)
	}
	return hash
}

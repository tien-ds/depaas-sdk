//go:build ios || dart
// +build ios dart

package main

import "C"
import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"math/big"
	"signer/key"
	"signer/mobile/dsshare"
	"signer/rpc"
	"signer/utils"
	"unsafe"
)

func main() {

}

//export ios_genKey
func ios_genKey() *C.char {
	return C.CString(genKey())
}

//export ios_genKeyFromSeed
func ios_genKeyFromSeed(seedHash *C.char) *C.char {
	return C.CString(genKeyFromSeed(C.GoString(seedHash)))
}

//export ios_privateKeyToAddr
func ios_privateKeyToAddr(base64PrivateKey *C.char) *C.char {
	return C.CString(key.PrivateKeyToAddr(C.GoString(base64PrivateKey)))
}

//export ios_hexToAddr
func ios_hexToAddr(sHex *C.char) *C.char {
	return C.CString(hexToAddr(C.GoString(sHex)))
}

//export ios_transfer
func ios_transfer(sJson *C.char) *C.char {
	return Call(func() interface{} {
		return transfer(C.GoString(sJson))
	})
}

//export ios_transfer1
func ios_transfer1(sJson *C.char) *C.char {
	return Call(func() interface{} {
		return transfer1(C.GoString(sJson))
	})
}

//export ios_getMapper
func ios_getMapper(sJson *C.char) *C.char {
	return Call(func() interface{} {
		return getMapper(C.GoString(sJson))
	})
}

//export ios_mapAccount
func ios_mapAccount(mainNetPrivateKey *C.char, base64PrivateKey *C.char) *C.char {
	return Call(func() interface{} {
		return mapAccount(C.GoString(mainNetPrivateKey), C.GoString(base64PrivateKey))
	})
}

//export ios_getBalance
func ios_getBalance(contractAddr *C.char, hexAddr *C.char) *C.char {
	return Call(func() interface{} {
		return getBalance(C.GoString(contractAddr), C.GoString(hexAddr))
	})
}

//export ios_tokenInfo
func ios_tokenInfo(contractAddr *C.char) *C.char {
	return Call(func() interface{} {
		return tokenInfo(C.GoString(contractAddr))
	})
}

//export ios_deposit
func ios_deposit(key *C.char, mainNetPrivateKey, tokenAmount *C.char) *C.char {
	return Call(func() interface{} {
		return deposit(C.GoString(key), C.GoString(mainNetPrivateKey), C.GoString(tokenAmount))
	})
}

//export ios_withdrawal
func ios_withdrawal(key *C.char, base64PrivateKey *C.char, mainNetPrivateKey *C.char, tokenAmount *C.char) *C.char {
	return Call(func() interface{} {
		return withdrawal(C.GoString(key), C.GoString(base64PrivateKey), C.GoString(mainNetPrivateKey), C.GoString(tokenAmount))
	})
}

//export ios_newContracts
func ios_newContracts(sMainNetCoin *C.char, sLoomCoin *C.char) *C.char {
	return Call(func() interface{} {
		return newContracts(C.GoString(sMainNetCoin), C.GoString(sLoomCoin))
	})
}

//export ios_eth_getTokenBalance
func ios_eth_getTokenBalance(url *C.char, contractAddr *C.char, addr *C.char) *C.char {
	return Call(func() interface{} {
		return rpc.GetTokenBalance(C.GoString(url), C.GoString(contractAddr), C.GoString(addr))
	})
}

//export ios_eth_getETHBalance
func ios_eth_getETHBalance(url *C.char, addr *C.char) *C.char {
	return Call(func() interface{} {
		return rpc.GetETHBalance(C.GoString(url), C.GoString(addr))
	})
}

//export ios_eth_allowance
func ios_eth_allowance(url *C.char, contractAddr *C.char, addr *C.char) *C.char {
	return Call(func() interface{} {
		return rpc.Allowance(C.GoString(url), C.GoString(contractAddr), C.GoString(addr))
	})
}

//export ios_eth_swapApprove
func ios_eth_swapApprove(contractAddr *C.char, mainNetPrivateKey *C.char) *C.char {
	return Call(func() interface{} {
		return rpc.SwapApprove(C.GoString(contractAddr), C.GoString(mainNetPrivateKey))
	})
}

//export ios_eth_getContractTokenInfo
func ios_eth_getContractTokenInfo(url *C.char, contract *C.char) *C.char {
	return Call(func() interface{} {
		return rpc.GetContractTokenInfoString(C.GoString(url), C.GoString(contract))
	})
}

//export ios_gateway_init
func ios_gateway_init() *C.char {
	return Call(func() interface{} {
		GateWayInit()
		return ""
	})
}

//export ios_ds_init
func ios_ds_init() *C.char {
	return Call(func() interface{} {
		DSInit()
		return ""
	})
}

//export ios_eth_transfer
func ios_eth_transfer(url *C.char, json *C.char) *C.char {
	return Call(func() interface{} {
		return ethTransfer(C.GoString(url), C.GoString(json))
	})
}

//export ios_eth_token_transfer
func ios_eth_token_transfer(contractAddr *C.char, mainNetPrivateKey *C.char, to *C.char, amount *C.char) *C.char {
	a, b := big.NewInt(0).SetString(C.GoString(amount), 10)
	utils.Require(b == true, "")
	return Call(func() interface{} {
		return rpc.EthTokenTransfer(C.GoString(contractAddr), C.GoString(mainNetPrivateKey), C.GoString(to), a)
	})
}

//export ios_commitcalltx
func ios_commitcalltx(contract *C.char, privateKey *C.char, input unsafe.Pointer, len C.int) *C.char {
	return Call(func() interface{} {
		return CommitCallTx(C.GoString(contract), C.GoString(privateKey), C.GoBytes(input, len))
	})
}

//export ios_set_contract
func ios_set_contract(contract *C.char) *C.char {
	return Call(func() interface{} {
		dsshare.SetDSContract(C.GoString(contract))
		return ""
	})

}

//export ios_ds_addStake
func ios_ds_addStake(mainPrivateKey, sender, amount *C.char) *C.char {
	return Call(func() interface{} {
		return dsshare.AddStake(C.GoString(mainPrivateKey), C.GoString(sender), C.GoString(amount))
	})
}

//export ios_ds_stakeOf
func ios_ds_stakeOf(account *C.char) *C.char {
	return Call(func() interface{} {
		return dsshare.DsStakeOf(C.GoString(account))
	})
}

//export ios_ds_coin_transfer
func ios_ds_coin_transfer(base64PrivateKey *C.char, to *C.char, amount *C.char) *C.char {
	return Call(func() interface{} {
		return DSCoinTransfer(C.GoString(base64PrivateKey), C.GoString(to), C.GoString(amount))
	})
}

//export ios_ds_coin_balance
func ios_ds_coin_balance(account *C.char) *C.char {
	return Call(func() interface{} {
		return DSCoinBalance(C.GoString(account))
	})
}

//export ios_sign
func ios_sign(rawHash *C.char, privateHash *C.char) *C.char {
	dataBytes := C.GoString(rawHash)
	fmt.Println("raw ", dataBytes)
	fmt.Println("pri ", dataBytes)
	priBytes := C.GoString(privateHash)
	pri, err := crypto.HexToECDSA(priBytes)
	if err != nil {
		panic(err)
	}
	//fmt.Println(crypto.PubkeyToAddress(pri.PublicKey).String())
	//Output: 0xa0D81025C4314f4E692cF71f9072955D46Efa8A5
	bytes, err := hex.DecodeString(dataBytes)
	if err != nil {
		panic(err)
	}
	sig, err := crypto.Sign(bytes, pri)
	if err != nil {
		panic(err)
	}
	return C.CString(hex.EncodeToString(sig))
}

func Call(callable func() interface{}) *C.char {
	return C.CString(toResult(func(stat *int, serr *string) interface{} {
		defer func() {
			if r := recover(); r != nil {
				*stat = 1
				if e, b := r.(error); b {
					*serr = e.Error()
				}
				if s, b := r.(string); b {
					*serr = s
				}
			}
		}()
		return callable()
	}))
}

func toResult(result func(stat *int, serr *string) interface{}) string {
	state := 0
	serr := ""
	k := result(&state, &serr)
	v := map[string]interface{}{
		"state":  state,
		"result": k,
		"err":    serr,
	}
	s, err := json.Marshal(v)
	if err != nil {
		logrus.Error(err)
	}
	return string(s)
}

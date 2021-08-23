package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"reflect"
	"strings"
)

func ToJson(f interface{}) string {
	bytes, err := json.Marshal(f)
	Require(err == nil, err)
	return string(bytes)
}

func Require(require bool, msg interface{}) {
	if !require {
		if e, b := msg.(error); b {
			panic(e)
		}
		kind := reflect.TypeOf(msg).Kind()
		if kind == reflect.String {
			panic(fmt.Errorf(msg.(string)))
		}
		if kind == reflect.Func {
			v := reflect.ValueOf(msg).Call(nil)
			if len(v) > 0 && v[0].Type().Kind() == reflect.String {
				fmt.Println(v[0].Interface().(string))
			}
		}
	}
}

func CosAddr(hex string) string {
	addr := "0x"
	if strings.HasPrefix(hex, "0x") || strings.HasPrefix(hex, "0X") || len(hex) == 40 {
		lAddr, err := loom.LocalAddressFromHexString(hex)
		Require(err == nil, err)
		addr += lAddr.Hex()
	}
	if len(hex) == 44 {
		bytes, err := base64.StdEncoding.DecodeString(hex)
		Require(err == nil, err)
		localAddress := loom.LocalAddressFromPublicKey(bytes)
		addr += localAddress.Hex()
	}
	if len(hex) == 28 {
		bytes, err := base64.StdEncoding.DecodeString(hex)
		Require(err == nil, err)
		localAddress := loom.LocalAddress(bytes)
		addr += localAddress.Hex()
	}
	if addr == "0x" {
		panic(fmt.Sprintf("%s cos not support addr", hex))
	}
	return addr
}

func ToAddrWithSigner(signer auth.Signer) string {
	return loom.LocalAddressFromPublicKey(signer.PublicKey()).String()
}

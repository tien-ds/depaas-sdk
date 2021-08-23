//go:build java
// +build java

package main

import (
	"gitee.com/aifuturewell/gojni/java"
	"signer/key"
	"signer/mobile/dsshare"
	"signer/rpc"
	"signer/utils"
)

func main() {

}

func init() {

	java.OnMainLoad(func(reg java.Register) {
		BaseJava(reg)
		DSErc20Java(reg)
		EthJava(reg)
		// OracleJava(reg)
		DEpaasJava(reg)
	})

}

func BaseJava(reg java.Register) {

	reg.WithClass("com.cos.Signer").
		BindNative("initClient", "void()", rpc.InitClient).
		BindNative("hexToAddr", "java.lang.String(java.lang.String)", hexToAddr).
		BindNative("addrToHex", "java.lang.String(java.lang.String)", AddrToHex).
		BindNative("privateKeyToAddr", "java.lang.String(java.lang.String)", key.PrivateKeyToAddr).
		BindNative("genKey", "java.lang.String()", genKey).
		BindNative("genKeyFromSeed", "java.lang.String(java.lang.String)", genKeyFromSeed).
		BindNative("getBalance", "java.lang.String(java.lang.String,java.lang.String)", getBalance).
		BindNative("transfer", "java.lang.String(java.lang.String)", transfer).
		BindNative("transfer1", "java.lang.String(java.lang.String)", transfer1).
		BindNative("getMapper", "java.lang.String(java.lang.String)", getMapper).
		BindNative("mapAccount", "java.lang.String(java.lang.String,java.lang.String)", mapAccount).
		BindNative("tokenInfo", "java.lang.String(java.lang.String)", tokenInfo).
		BindNative("sign", "java.lang.String(java.lang.String,byte[])", sign).
		BindNative("commitCallTx", "java.lang.String(java.lang.String,java.lang.String,byte[])", CommitCallTx).
		Done()
}

func OracleJava(reg java.Register) {

	reg.WithClass("com.cos.Oracle").
		BindNative("gatewayInit", "void()", GateWayInit).
		BindNative("deposit", "java.lang.String(java.lang.String,java.lang.String,java.lang.String)", deposit).
		BindNative("withdrawal", "java.lang.String(java.lang.String,java.lang.String,java.lang.String,java.lang.String)", withdrawal).
		BindNative("newContracts", "java.lang.String(java.lang.String,java.lang.String)", newContracts).
		Done()

}

// DSErc20Java ds miner
func DSErc20Java(reg java.Register) {

	//mainPrivateKey, sender, recipient, amount string
	reg.WithClass("com.ds.DSErc20").
		BindNative("setContract", "void(java.lang.String)", dsshare.SetDSContract).
		BindNative("stackOf", "java.lang.String(java.lang.String)", dsshare.DsStakeOf).
		BindNative("getMinerAddr", "java.lang.String(java.lang.String)", dsshare.DsGetMinerAddr).
		BindNative("addStack", "java.lang.String(java.lang.String,java.lang.String,java.lang.String)", dsshare.AddStake).
		BindNative("penalty", "java.lang.String(java.lang.String,java.lang.String,java.lang.String)", dsshare.DsPenalty).
		BindNative("exchange", "java.lang.String(java.lang.String,java.lang.String,java.lang.String,java.lang.String)", dsshare.DsExchange).
		BindNative("owner", "java.lang.String()", dsshare.DsOwner).
		Done()

}

func DEpaasJava(reg java.Register) {
	reg.WithClass("com.ds.Coin").
		BindNative("mint", "java.lang.String(java.lang.String,java.lang.String,java.lang.String)", DSCoinMint).
		BindNative("balance", "java.lang.String(java.lang.String)", DSCoinBalance).
		BindNative("transfer", "java.lang.String(java.lang.String,java.lang.String,java.lang.String)", DSCoinTransfer).
		Done()
}

func EthJava(reg java.Register) {

	reg.WithClass("com.cos.ETH").
		BindNative("getTokenBalance", "java.lang.String(java.lang.String,java.lang.String,java.lang.String)", rpc.GetTokenBalance).
		BindNative("getContractTokenInfo", "java.lang.String(java.lang.String,java.lang.String)", func(url string, contract string) string {
			return utils.ToJson(rpc.GetContractTokenInfo(url, contract))
		}).
		BindNative("getEthBalance", "java.lang.String(java.lang.String,java.lang.String)", rpc.GetETHBalance).
		BindNative("allowance", "java.lang.String(java.lang.String,java.lang.String,java.lang.String)", rpc.Allowance).
		BindNative("swapApprove", "java.lang.String(java.lang.String,java.lang.String)", rpc.SwapApprove).
		Done()

}

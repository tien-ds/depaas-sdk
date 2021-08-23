package rpc

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/loomnetwork/go-loom/client"
	"math"
	"math/big"
	"signer/glob"
	"signer/token/mainnet"
	"signer/utils"
	"strings"
)

func Init(url string, contractAddr string) *mainnet.MainnetERC20MintableContract {
	client, err := ethclient.Dial(url)
	utils.Require(err == nil, err)
	contract, err := mainnet.ConnectToMainnetERC20MintableContract(client, contractAddr)
	return contract
}

func GetTokenBalance(url string, contractAddr string, addr string) string {
	b, e := Init(url, contractAddr).BalanceOf(&client.Identity{
		MainnetAddr: common.HexToAddress(addr),
	})
	utils.Require(e == nil, e)
	return b.String()
}

func Allowance(url string, contractAddr string, addr string) string {
	bls, err := Init(url, contractAddr).Allowance(common.HexToAddress(addr), common.HexToAddress(glob.GATEWAY))
	utils.Require(err == nil, err)
	return bls.String()
}

func SwapApprove(contractAddr string, mainNetPrivateKey string) string {
	mainNetPrivKey, e := crypto.HexToECDSA(strings.TrimPrefix(mainNetPrivateKey, "0x"))
	utils.Require(e == nil, e)
	contract, err := mainnet.NewSampleERC20MintableToken(common.HexToAddress(contractAddr), ethClient)
	utils.Require(err == nil, err)
	tx, err := contract.Approve(client.DefaultTransactOptsForIdentity(&client.Identity{
		MainnetPrivKey: mainNetPrivKey,
	}), common.HexToAddress(glob.GATEWAY), big.NewInt(math.MinInt32))
	utils.Require(err == nil, err)
	return tx.Hash().String()
}

func EthTokenTransfer(contractAddr string, mainNetPrivateKey string, to string, amount *big.Int) string {
	mainNetPrivKey, e := crypto.HexToECDSA(strings.TrimPrefix(mainNetPrivateKey, "0x"))
	utils.Require(e == nil, e)
	contract, err := mainnet.NewSampleERC20MintableToken(common.HexToAddress(contractAddr), ethClient)
	utils.Require(err == nil, err)
	tx, err := contract.Transfer(client.DefaultTransactOptsForIdentity(&client.Identity{
		MainnetPrivKey: mainNetPrivKey,
	}), common.HexToAddress(to), amount)
	utils.Require(err == nil, err)
	return tx.Hash().String()
}

func EthTransfer(url string, to *common.Address, amount *big.Int, chainId *big.Int, privateKey string) string {
	client, err := ethclient.Dial(url)
	utils.Require(err == nil, err)
	priKey, err := crypto.HexToECDSA(privateKey)
	utils.Require(err == nil, err)
	from := crypto.PubkeyToAddress(priKey.PublicKey)
	context := context.Background()
	nonce, err := client.PendingNonceAt(context, from)
	utils.Require(err == nil, err)
	gasPrice, err := client.SuggestGasPrice(context)
	utils.Require(err == nil, err)
	msg := ethereum.CallMsg{From: from, To: to, Value: amount}
	gasLimit, err := client.EstimateGas(context, msg)
	utils.Require(err == nil, err)
	tx := types.NewTransaction(nonce, *to, amount, gasLimit, gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), priKey)
	utils.Require(err == nil, err)
	utils.Require(client.SendTransaction(context, signedTx) == nil, "")
	return signedTx.Hash().String()
}

func GetETHBalance(url string, addr string) string {
	client, err := ethclient.Dial(url)
	utils.Require(err == nil, err)
	big, err := client.BalanceAt(context.Background(), common.HexToAddress(addr), nil)
	utils.Require(err == nil, err)
	return big.String()
}

func GetContractTokenInfoString(url string, contract string) string {
	return utils.ToJson(GetContractTokenInfo(url, contract))
}

func GetContractTokenInfo(url string, contractAddr string) TokenInfo {
	client, err := ethclient.Dial(url)
	utils.Require(err == nil, err)
	contract, err := mainnet.NewSampleERC20MintableToken(common.HexToAddress(contractAddr), client)
	utils.Require(err == nil, err)
	name, err := contract.Name(nil)
	utils.Require(err == nil, err)
	decimals, err := contract.Decimals(nil)
	utils.Require(err == nil, err)
	symbol, err := contract.Symbol(nil)
	utils.Require(err == nil, err)
	return TokenInfo{
		Name:     name,
		Decimals: decimals,
		Symbol:   symbol,
	}
}

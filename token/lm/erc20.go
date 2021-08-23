package lm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
	"math/big"
	"signer/utils"
)

var ERC20ABI = `[{"inputs":[{"internalType":"string","name":"name","type":"string"},{"internalType":"string","name":"symbol","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"}]`

type DAppChainEVMContract struct {
	*ContractConnect
}

func NewDAppChainEVMContract(contract string, dClient *client.DAppChainRPCClient) *DAppChainEVMContract {
	la, _ := loom.LocalAddressFromHexString(contract)
	adc := NewDAppChainContract(dClient, ERC20ABI, loom.Address{
		ChainID: dClient.GetChainID(),
		Local:   la,
	})
	return &DAppChainEVMContract{adc}
}

func (c *DAppChainEVMContract) BalanceOf(hexAddr string) (*big.Int, error) {
	ownerAddr := common.HexToAddress(utils.CosAddr(hexAddr))
	var result *big.Int
	if err := c.StaticCallEVM("balanceOf", &result, ownerAddr); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *DAppChainEVMContract) Name() string {
	var name string
	if err := c.StaticCallEVM("name", &name); err != nil {
		panic(err)
	}
	return name
}

func (c *DAppChainEVMContract) Decimals() uint8 {
	var num uint8
	if err := c.StaticCallEVM("decimals", &num); err != nil {
		return 0
	}
	return num
}

func (c *DAppChainEVMContract) Symbol() string {
	var name string
	if err := c.StaticCallEVM("symbol", &name); err != nil {
		return ""
	}
	return name
}

// TotalSupply totalSupply
func (c *DAppChainEVMContract) TotalSupply() (*big.Int, error) {
	var result *big.Int
	err := c.StaticCallEVM("totalSupply", &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Approve grants authorization to an entity to transfer the given tokens at a later time
func (c *DAppChainEVMContract) Approve(signer auth.Signer, to string, amount *big.Int) (string, error) {
	toAddr := common.HexToAddress(utils.CosAddr(to))
	return c.CallEVM("approve", signer, toAddr, amount)
}

func (c *DAppChainEVMContract) Transfer(signer auth.Signer, to string, amount *big.Int) (string, error) {
	toAddr := common.HexToAddress(utils.CosAddr(to))
	return c.CallEVM("transfer", signer, toAddr, amount)
}

func (c *DAppChainEVMContract) MintTo(signer auth.Signer, to string, amount *big.Int) (string, error) {
	toAddr := common.HexToAddress(utils.CosAddr(to))
	return c.CallEVM("mintTo", signer, toAddr, amount)
}

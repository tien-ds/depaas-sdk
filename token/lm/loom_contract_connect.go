package lm

import (
	"encoding/base64"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
	"signer/utils"
	"strings"
)

type ContractConnect struct {
	Contract    *client.EvmContract
	ContractABI *abi.ABI
	ChainID     string
	Address     loom.Address
}

func (c *ContractConnect) StaticCallEVM(method string, result interface{}, params ...interface{}) error {
	input, err := c.ContractABI.Pack(method, params...)
	if err != nil {
		return err
	}
	output, err := c.Contract.StaticCall(input, c.Address)
	if err != nil {
		return err
	}
	return c.ContractABI.Unpack(result, method, output)
}

func (c *ContractConnect) CallEVM(method string, signer auth.Signer, params ...interface{}) (string, error) {
	input, err := c.ContractABI.Pack(method, params...)
	if err != nil {
		return "", err
	}
	bytes, err := c.Contract.Call(input, signer)
	return base64.StdEncoding.EncodeToString(bytes), err
}

func NewDAppChainContract(rpcClient *client.DAppChainRPCClient, ContractABI string, contractAddr loom.Address) *ContractConnect {
	if rpcClient == nil {
		panic("client is nil")
	}
	abi, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		panic(err)
	}
	contract := client.NewEvmContract(rpcClient, contractAddr.Local)
	utils.Require(contract != nil, fmt.Sprintf("add %s is not contract", contractAddr.String()))
	return &ContractConnect{
		Contract:    contract,
		ContractABI: &abi,
		ChainID:     rpcClient.GetChainID(),
		Address:     contractAddr,
	}
}

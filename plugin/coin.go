package plugin

import (
	"encoding/base64"
	"math/big"
	"signer/utils"

	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/builtin/types/coin"
	"github.com/loomnetwork/go-loom/client"
	"github.com/loomnetwork/go-loom/types"
)

type CoinPlugin struct {
	contract *client.Contract
}

func (c *CoinPlugin) BalanceOf(hexAddr string) (*big.Int, error) {

	addr, err := loom.LocalAddressFromHexString(utils.CosAddr(hexAddr))
	if err != nil {
		return nil, err
	}
	param := coin.BalanceOfRequest{
		Owner: &types.Address{
			ChainId: "default",
			Local:   addr,
		},
	}
	var resp coin.BalanceOfResponse
	c.contract.StaticCall("BalanceOf", &param, loom.Address{
		ChainID: "default",
		Local:   addr,
	}, &resp)
	return resp.Balance.Value.Int, nil
}

func (c *CoinPlugin) Name() string {
	panic("Name not Name")
}

func (c *CoinPlugin) Decimals() uint8 {
	panic("Name not Decimals")
}

func (c *CoinPlugin) Symbol() string {
	panic("Name not Symbol")
}

func (c *CoinPlugin) Approve(signer auth.Signer, to string, amount *big.Int) (string, error) {
	panic("Name not Approve")
}

func (c *CoinPlugin) Mint(signer auth.Signer, to string, amount *big.Int) (string, error) {
	localAddr, err := loom.LocalAddressFromHexString(utils.CosAddr(to))
	if err != nil {
		return "", err
	}
	lAddr := loom.Address{
		ChainID: "default",
		Local:   localAddr,
	}
	transParam := &coin.TransferRequest{
		To: lAddr.MarshalPB(),
		Amount: &types.BigUInt{
			Value: *loom.NewBigUInt(amount),
		},
	}
	hash, err := c.contract.Call("Mint", transParam, signer, nil)
	if err != nil {
		return "", err
	}
	if byts, b := hash.([]uint8); b {
		return base64.StdEncoding.EncodeToString(byts), nil
	}
	return "", nil
}

func (c *CoinPlugin) Transfer(signer auth.Signer, to string, amount *big.Int) (string, error) {
	localAddr, err := loom.LocalAddressFromHexString(utils.CosAddr(to))
	if err != nil {
		return "", err
	}
	lAddr := loom.Address{
		ChainID: "default",
		Local:   localAddr,
	}
	transParam := &coin.TransferRequest{
		To: lAddr.MarshalPB(),
		Amount: &types.BigUInt{
			Value: *loom.NewBigUInt(amount),
		},
	}
	hash, err := c.contract.Call("Transfer", transParam, signer, nil)
	if err != nil {
		return "", err
	}
	if byts, b := hash.([]uint8); b {
		return base64.StdEncoding.EncodeToString(byts), nil
	}
	return "", nil
}

func NewCoinPlugin(clientDapp *client.DAppChainRPCClient, contractAddr loom.LocalAddress) *CoinPlugin {
	return &CoinPlugin{
		contract: client.NewContract(clientDapp, contractAddr),
	}
}

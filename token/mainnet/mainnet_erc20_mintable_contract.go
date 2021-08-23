package mainnet

import (
	"context"
	"github.com/loomnetwork/go-loom/client"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type MainnetERC20MintableContract struct {
	contract  *SampleERC20MintableToken
	ethClient *ethclient.Client

	TxTimeout time.Duration
	Address   common.Address
	TxHash    string
}

func (c *MainnetERC20MintableContract) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	bal, err := c.contract.Allowance(nil, owner, spender)
	if err != nil {
		return nil, err
	}
	return bal, nil
}

func (c *MainnetERC20MintableContract) BalanceOf(caller *client.Identity) (*big.Int, error) {
	bal, err := c.contract.BalanceOf(nil, caller.MainnetAddr)
	if err != nil {
		return nil, err
	}
	return bal, nil
}

func (c *MainnetERC20MintableContract) TransferFrom(to *client.Identity, from *client.Identity, amount *big.Int) error {
	tx, err := c.contract.TransferFrom(client.DefaultTransactOptsForIdentity(from), from.MainnetAddr, to.MainnetAddr, amount)
	if err != nil {
		return err
	}
	return client.WaitForTxConfirmation(context.TODO(), c.ethClient, tx, c.TxTimeout)
}

func (c *MainnetERC20MintableContract) Approve(from *client.Identity, to common.Address, amount *big.Int) error {
	tx, err := c.contract.Approve(client.DefaultTransactOptsForIdentity(from), to, amount)
	if err != nil {
		return err
	}
	return client.WaitForTxConfirmation(context.TODO(), c.ethClient, tx, c.TxTimeout)
}

func (c *MainnetERC20MintableContract) Mint(from *client.Identity, to common.Address, amount *big.Int) error {
	tx, err := c.contract.Mint(client.DefaultTransactOptsForIdentity(from), to, amount)
	if err != nil {
		return err
	}
	return client.WaitForTxConfirmation(context.TODO(), c.ethClient, tx, c.TxTimeout)
}

func (c *MainnetERC20MintableContract) MintTo(from *client.Identity, to common.Address, amount *big.Int) error {
	tx, err := c.contract.MintTo(client.DefaultTransactOptsForIdentity(from), to, amount)
	if err != nil {
		return err
	}
	return client.WaitForTxConfirmation(context.TODO(), c.ethClient, tx, c.TxTimeout)
}

func (c *MainnetERC20MintableContract) Transfer(to *client.Identity, from *client.Identity, amount *big.Int) error {
	tx, err := c.contract.Transfer(client.DefaultTransactOptsForIdentity(from), to.MainnetAddr, amount)
	if err != nil {
		return err
	}
	return client.WaitForTxConfirmation(context.TODO(), c.ethClient, tx, c.TxTimeout)
}

func ConnectToMainnetERC20MintableContract(ethClient *ethclient.Client, address string) (*MainnetERC20MintableContract, error) {
	contractAddr := common.HexToAddress(address)
	contract, err := NewSampleERC20MintableToken(contractAddr, ethClient)
	if err != nil {
		return nil, err
	}
	return &MainnetERC20MintableContract{
		contract:  contract,
		ethClient: ethClient,
		Address:   contractAddr,
	}, nil
}

package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"signer/glob"
	"signer/key"
	"signer/token/lm"

	"signer/token/mainnet"
	"signer/utils"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/client"
	gw "github.com/loomnetwork/go-loom/client/gateway_v2"
	vmc "github.com/loomnetwork/go-loom/client/validator_manager"
	"github.com/sirupsen/logrus"
)

type Swap struct {
	ethClient          *ethclient.Client
	mainNetCoin        *mainnet.MainnetERC20MintableContract
	loomCoin           *lm.DAppChainEVMContract
	validatorsManager  *vmc.MainnetVMCClient
	dAppChainRPCClient *client.DAppChainRPCClient
	mainnetGateway     *gw.MainnetGatewayClient
	dappChainId        string
}

func NewSwap(ethClient *ethclient.Client, dAppChainRPCClient *client.DAppChainRPCClient) *Swap {

	mainnetGateway, e := gw.ConnectToMainnetGateway(ethClient, glob.GATEWAY)
	utils.Require(e == nil, e)
	validatorsManager, e := vmc.ConnectToMainnetVMCClient(ethClient, glob.VALIDATORMANAGER)
	utils.Require(e == nil, e)
	return &Swap{
		ethClient:          ethClient,
		validatorsManager:  validatorsManager,
		dAppChainRPCClient: dAppChainRPCClient,
		mainnetGateway:     mainnetGateway,
	}
}

//0x0a37b2fdcc10c80da61804eb7da01b5ad3c66ecc,0xc1b174f1bc70172c911df8a01d0e0d98129c7517
func (s *Swap) SetContracts(sMainNetCoin string, sLoomCoin string) {
	dappChainId := s.dAppChainRPCClient.GetChainID()
	s.dappChainId = dappChainId
	loomCoin, e := loom.LocalAddressFromHexString(sLoomCoin)
	utils.Require(e == nil, e)
	loomERC20 := &lm.DAppChainEVMContract{ContractConnect: lm.NewDAppChainContract(s.dAppChainRPCClient, lm.ERC20ABI, loom.Address{
		Local:   loomCoin,
		ChainID: dappChainId,
	})}
	s.loomCoin = loomERC20
	mainnetCoin, e := mainnet.ConnectToMainnetERC20MintableContract(s.ethClient, sMainNetCoin)
	utils.Require(e == nil, e)
	s.mainNetCoin = mainnetCoin
}

func (s *Swap) Deposit(mainNetPrivateKey string, tokenAmount *big.Int) error {
	start := time.Now()
	mainNetPrivKey, e := crypto.HexToECDSA(strings.TrimPrefix(mainNetPrivateKey, "0x"))
	utils.Require(e == nil, e)
	alice := &client.Identity{
		MainnetPrivKey: mainNetPrivKey,
	}
	has, err := s.mainNetCoin.Allowance(crypto.PubkeyToAddress(mainNetPrivKey.PublicKey), s.mainnetGateway.Address)
	utils.Require(err == nil, err)
	if has.Cmp(tokenAmount) <= 0 {
		return utils.NewRPCError(errors.New("not Approve"), 0)
	}
	tx, err := s.mainnetGateway.Contract().DepositERC20(client.DefaultTransactOptsForIdentity(alice), tokenAmount, s.mainNetCoin.Address)
	if err != nil {
		return err
	}
	logrus.Infof("Deposit use %s tx %s", time.Since(start), tx.Hash().String())
	json, err := tx.MarshalJSON()
	utils.Require(err == nil, err)
	return utils.NewSuccess(string(json))
}

func (s *Swap) ResumeWithdrawal(base64PrivateKey string, mainNetPrivateKey string) error {
	accountSigner := key.NewSigner(base64PrivateKey)
	publicKey := loom.LocalAddressFromPublicKey(accountSigner.PublicKey())
	mainnetPrivKey, e := crypto.HexToECDSA(strings.TrimPrefix(mainNetPrivateKey, "0x"))
	utils.Require(e == nil, e)
	alice := &client.Identity{
		LoomSigner: accountSigner,
		LoomAddr: loom.Address{
			Local:   publicKey,
			ChainID: s.dappChainId,
		},
		MainnetAddr:    crypto.PubkeyToAddress(mainnetPrivKey.PublicKey),
		MainnetPrivKey: mainnetPrivKey,
	}
	dappChainGateway, e := gw.ConnectToDAppChainGateway(s.dAppChainRPCClient, glob.EVENTSURI)
	utils.Require(e == nil, e)
	receipt, e := dappChainGateway.WithdrawalReceipt(alice)
	if receipt != nil {
		fmt.Println(common.HexToAddress(receipt.TokenContract.Local.Hex()).Hex())
		fmt.Println(s.mainNetCoin.Address.String())
		fmt.Printf("Found pending withdrawal of %s coins\n", receipt.TokenAmount.Value.String())
		hex := common.HexToAddress(receipt.TokenContract.Local.Hex()).Hex()
		if s.mainNetCoin.Address.String() == hex {
			fmt.Println("eth")
			validators, err := s.validatorsManager.GetValidators()
			err = s.mainnetGateway.WithdrawERC20(alice, receipt.TokenAmount.Value.Int, s.mainNetCoin.Address, receipt.OracleSignature, validators)
			if err == nil {
				logrus.Info("Withdrawal success")
			} else {
				return err
			}
		}
		if s.loomCoin.Address.Local.Hex() == hex {
			fmt.Println("loom")
		}
		if s.mainnetGateway.Address.String() == hex {
			fmt.Println("gateway")
		}
	}
	return nil
}

// 3 cos balance 0
// 2 eth balance 0
// 4 has pending Withdrawal
// 1 unknown error
//e5c1c6c2e7b7af7dce8a5fbeadae0bfb83834a7c3c8d0696a3ee5c6edd354f25
//Go0cxZ4WwmU7nXmRwKVnDnED44LWxbijH7I+Jp8iY8CDIHpnz+xVzkAKxshvS8quo/ZAKMeouyuJtb5X2DSc+g==
func (s *Swap) Withdrawal(base64PrivateKey string, mainNetPrivateKey string, tokenAmount *big.Int) error {
	mapperAddr, e := s.dAppChainRPCClient.Resolve("gateway")
	utils.Require(e == nil, e)
	dappchainGateway, e := gw.ConnectToDAppChainGateway(s.dAppChainRPCClient, glob.EVENTSURI)
	utils.Require(e == nil, e)
	mainnetPrivKey, e := crypto.HexToECDSA(strings.TrimPrefix(mainNetPrivateKey, "0x"))
	utils.Require(e == nil, e)

	accountSigner := key.NewSigner(base64PrivateKey)
	publicKey := loom.LocalAddressFromPublicKey(accountSigner.PublicKey())
	alice := &client.Identity{
		LoomSigner: accountSigner,
		LoomAddr: loom.Address{
			Local:   publicKey,
			ChainID: s.dappChainId,
		},
		MainnetAddr:    crypto.PubkeyToAddress(mainnetPrivKey.PublicKey),
		MainnetPrivKey: mainnetPrivKey,
	}

	addr := loom.LocalAddressFromPublicKey(accountSigner.PublicKey()).String()

	balance, e := s.loomCoin.BalanceOf(addr)
	if e != nil {
		return e
	}
	logrus.Infof("cos %s is %s", addr, balance)
	if balance == nil || balance.Int64() == 0 {
		logrus.Warnf("%s balance is 0", addr)
		return utils.NewRPCError(fmt.Errorf("cos %s is %s", addr, balance), 3)
	}

	ethBalance, e := s.ethClient.BalanceAt(context.Background(), alice.MainnetAddr, nil)
	utils.Require(e == nil, e)
	logrus.Infof("eth %s balance is %s", alice.MainnetAddr.String(), ethBalance)
	if ethBalance == nil || ethBalance.Cmp(big.NewInt(0)) == 0 {
		return utils.NewRPCError(fmt.Errorf("eth %s balance is %s", alice.MainnetAddr.String(), ethBalance), 2)
	}

	receipt, e := dappchainGateway.WithdrawalReceipt(alice)
	if receipt != nil {
		fmt.Println(receipt.TokenContract.Local.Hex())
		fmt.Println(s.mainNetCoin.Address.Hex())
		return utils.NewRPCError(fmt.Errorf("found pending withdrawal of %s coins", receipt.TokenAmount.Value.String()), 4)
	}

	hex, err := s.loomCoin.Approve(accountSigner, mapperAddr.Local.String(), tokenAmount)
	if err != nil {
		return err
	}
	logrus.Infof("loomCoin.Approve %s", hex)
	for {
		wr, e := dappchainGateway.WithdrawalReceipt(alice)
		if e != nil {
			logrus.Error(e)
			return e
		}
		if wr == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}

	for i := 0; i < 5; i++ {
		err := dappchainGateway.WithdrawERC20(alice, tokenAmount, s.loomCoin.Address)
		if err != nil {
			if strings.Contains(err.Error(), "TG003") {
				logrus.Warn("dappchainGateway.WithdrawERC20 sleep 5s")
				time.Sleep(5 * time.Second)
			} else {
				logrus.Error(err)
				return err
			}
		} else {
			break
		}
	}

	wr, err := dappchainGateway.WithdrawalReceipt(alice)
	if wr == nil {
		logrus.Warn("DAppChain Mapper should've cleared out Alice's pending withdrawal")
		return err
	}
	for {
		wr, err = dappchainGateway.WithdrawalReceipt(alice)
		if err != nil {
			logrus.Error(err)
			return err
		}
		if wr != nil && wr.OracleSignature != nil {
			break
		}
		logrus.Warn("dappchainGateway.WithdrawalReceipt(alice) sleep 5s")
		time.Sleep(5 * time.Second)
	}
	validators, err := s.validatorsManager.GetValidators()
	if err != nil {
		return utils.NewRPCError(err, 0)
	}
	c := s.mainnetGateway.Contract()

	hash := withdrawalHash(s.mainnetGateway, alice.MainnetAddr, s.mainNetCoin.Address, big.NewInt(0), tokenAmount)
	v, r, s1, valIndexes, err := client.ParseSigs(wr.OracleSignature, hash, validators)
	if err != nil {
		return utils.NewRPCError(err, 0)
	}
	tx, err := c.WithdrawERC20(client.DefaultTransactOptsForIdentity(alice), tokenAmount, s.mainNetCoin.Address, valIndexes, v, r, s1)
	if err != nil {
		return utils.NewRPCError(err, 0)
	}
	trx, err := tx.MarshalJSON()
	utils.Require(err == nil, err)
	return utils.NewSuccess(string(trx))
}

func withdrawalHash(c *gw.MainnetGatewayClient, withdrawer common.Address, tokenAddr common.Address, tokenId *big.Int, amount *big.Int) []byte {
	nonce, err := c.Nonces(withdrawer)
	if err != nil {
		return nil
	}
	hash := client.WithdrawalHash(withdrawer, tokenAddr, c.Address, 1, tokenId, amount, nonce, true)
	return client.ToEthereumSignedMessage(hash)
}

//08635939

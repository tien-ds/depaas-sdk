package gateway

import (
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/builtin/types/address_mapper"
	"github.com/loomnetwork/go-loom/client"
	"github.com/loomnetwork/go-loom/common/evmcompat"
	ssha "github.com/miguelmota/go-solidity-sha3"
	"github.com/sirupsen/logrus"
	"signer/key"
	"signer/utils"
	"strings"
)

var encoder = base64.StdEncoding

type Mapper struct {
	cli            *client.DAppChainRPCClient
	mapperContract *client.Contract
}

func NewMapper(rpcClient *client.DAppChainRPCClient) *Mapper {
	mapperAddr, e := rpcClient.Resolve("addressmapper")
	utils.Require(e == nil, e)
	mapperContract := client.NewContract(rpcClient, mapperAddr.Local)

	return &Mapper{
		cli:            rpcClient,
		mapperContract: mapperContract,
	}
}

func (c *Mapper) GetMappedAccount(account loom.Address) (loom.Address, error) {
	req := &address_mapper.AddressMapperGetMappingRequest{
		From: account.MarshalPB(),
	}
	resp := &address_mapper.AddressMapperGetMappingResponse{}
	_, err := c.mapperContract.StaticCall("GetMapping", req, account, resp)
	if err != nil {
		return loom.Address{}, err
	}
	return loom.UnmarshalAddressPB(resp.To), nil
}

func (c *Mapper) MapAccounts(mainNetPrivateKey string, base64PrivateKey string) error {
	mainnetPrivKey, _ := crypto.HexToECDSA(strings.TrimPrefix(mainNetPrivateKey, "0x"))
	mainnetLocalAddr, err := loom.LocalAddressFromHexString(crypto.PubkeyToAddress(mainnetPrivKey.PublicKey).String())
	if err != nil {
		return err
	}
	from := loom.Address{
		ChainID: "eth",
		Local:   mainnetLocalAddr,
	}
	accountSigner := key.NewSigner(base64PrivateKey)
	to := loom.Address{
		ChainID: c.cli.GetChainID(),
		Local:   loom.LocalAddressFromPublicKey(accountSigner.PublicKey()),
	}

	mappedAccount, err := c.GetMappedAccount(from)
	if err == nil {
		if mappedAccount.Compare(to) != 0 {
			logrus.Warnf("Account %v is mapped to %v", from, mappedAccount)
			return fmt.Errorf("Account %v is mapped to %v", from, mappedAccount)
		}
		return fmt.Errorf("Account %v has mapped to %v", from, mappedAccount)
	}

	logrus.Infof("Mapping account %v to %v\n", from, to)

	sig, err := signIdentityMapping(from, to, mainnetPrivKey)
	if err != nil {
		return err
	}
	req := &address_mapper.AddressMapperAddIdentityMappingRequest{
		From:      from.MarshalPB(),
		To:        to.MarshalPB(),
		Signature: sig,
	}
	_, err = c.mapperContract.Call("AddIdentityMapping", req, accountSigner, nil)
	if err != nil {
		logrus.Error(err)
	}
	return err
}

func signIdentityMapping(from, to loom.Address, key *ecdsa.PrivateKey) ([]byte, error) {
	hash := ssha.SoliditySHA3(
		ssha.Address(common.BytesToAddress(from.Local)),
		ssha.Address(common.BytesToAddress(to.Local)),
	)
	sig, err := evmcompat.SoliditySign(hash, key)
	if err != nil {
		return nil, err
	}
	// Prefix the sig with a single byte indicating the sig type, in this case EIP712
	return append(make([]byte, 1, 66), sig...), nil
}

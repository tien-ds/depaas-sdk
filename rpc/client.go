package rpc

import (
	"context"
	"encoding/base64"
	"errors"
	"math/big"
	"net"
	"net/http"
	"os"
	"signer/gateway"
	"signer/glob"
	"signer/key"
	"signer/token/lm"
	"signer/utils"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/client"
)

var (
	dAppClient *client.DAppChainRPCClient
	swap       *gateway.Swap
	mapper     *gateway.Mapper
	ethClient  *ethclient.Client
)

func InitClient() {
	dAppClient = client.NewDAppChainRPCClientShareTransport(
		glob.ChainId,
		glob.WRITEURI,
		glob.READURI,
		&http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
				return net.DialTimeout(network, addr, time.Second*3)
			},
		},
	)
	mapper = gateway.NewMapper(dAppClient)
}

func SetEthClient(url *string) {
	var err error
	if url != nil {
		ethClient, err = ethclient.Dial(*url)
	} else {
		ethClient, err = ethclient.Dial(glob.ETHEREUM_NETWORK)
	}
	utils.Require(err == nil, err)
}

func GateWayInit() {
	SetEthClient(nil)
	os.Setenv("ETHEREUM_NETWORK", glob.ETHEREUM_NETWORK)
	swap = gateway.NewSwap(ethClient, dAppClient)
	//swap.SetContracts("0x0a37b2fdcc10c80da61804eb7da01b5ad3c66ecc", "0xc1b174f1bc70172c911df8a01d0e0d98129c7517")
}

func GetMainNetClient() *ethclient.Client {
	return ethClient
}

func GetDAppClient() *client.DAppChainRPCClient {
	return dAppClient
}

func GetAccountMP() *gateway.Mapper {
	return mapper
}

// GetSwap TODO
func GetSwap() *gateway.Swap {
	return swap
}

type Result struct {
	State  int         `json:"state"`
	Error  interface{} `json:"error"`
	Result interface{} `json:"result"`
}

func ToMessageResult(err error, message interface{}) Result {
	if e, b := err.(utils.RPCError); b {
		return Result{
			State: e.GetState(),
			Error: e.Error(),
		}
	}
	if err == nil {
		return Result{
			State:  1,
			Result: message,
		}
	} else {
		return Result{
			State: 0,
			Error: err.Error(),
		}
	}
}

func ToResult(err error) Result {
	if e, b := err.(utils.RPCError); b {
		return Result{
			State:  e.GetState(),
			Error:  e.Error(),
			Result: e.GetMessage(),
		}
	}
	return ToMessageResult(err, nil)
}

func AccountMapper(param struct {
	MainNetPrivateKey string `json:"mainNetPrivateKey"`
	Base64PrivateKey  string `json:"base64PrivateKey"`
}) Result {
	err := mapper.MapAccounts(param.MainNetPrivateKey, param.Base64PrivateKey)
	return ToResult(err)
}

func GetMapper(param struct {
	Type    int    `json:"type"`
	HexAddr string `json:"hexAddr"`
}) Result {
	t := ""
	if param.Type == 0 {
		t = "eth"
	} else {
		t = glob.ChainId
	}
	la, _ := loom.LocalAddressFromHexString(utils.CosAddr(param.HexAddr))
	mapAdds, e := mapper.GetMappedAccount(loom.Address{
		ChainID: t,
		Local:   la,
	})
	if param.Type == 0 {
		return ToMessageResult(e, map[string]string{"addr": base64.StdEncoding.EncodeToString(mapAdds.Local)})
	} else {
		return ToMessageResult(e, map[string]string{"addr": mapAdds.Local.String()})
	}
}

func Transfer(param struct {
	ContractAddr     string `json:"contractAddr"`
	To               string `json:"to"`
	Base64PrivateKey string `json:"base64PrivateKey"`
	Amount           string `json:"amount"`
}) Result {
	la, _ := loom.LocalAddressFromHexString(utils.CosAddr(param.ContractAddr))
	adc := lm.NewDAppChainContract(dAppClient, lm.ERC20ABI, loom.Address{
		ChainID: dAppClient.GetChainID(),
		Local:   la,
	})
	contr := lm.DAppChainEVMContract{ContractConnect: adc}
	am, b := big.NewInt(0).SetString(param.Amount, 10)
	if !b {
		return ToMessageResult(errors.New("big convert error"), "")
	}
	s, e := contr.Transfer(
		key.NewSigner(param.Base64PrivateKey),
		param.To,
		am,
	)
	return ToMessageResult(e, map[string]string{"hash": s})
}

func GetBalance(contractAddr string, hexAddr string) *big.Int {
	la, _ := loom.LocalAddressFromHexString(utils.CosAddr(contractAddr))
	adc := lm.NewDAppChainContract(dAppClient, lm.ERC20ABI, loom.Address{
		ChainID: dAppClient.GetChainID(),
		Local:   la,
	})
	contr := lm.DAppChainEVMContract{ContractConnect: adc}
	//9995800000000000000000
	amount, err := contr.BalanceOf(hexAddr)
	utils.Require(err == nil, err)
	return amount
}

type TokenInfo struct {
	Name     string `json:"name"`
	Decimals uint8  `json:"decimals"`
	Symbol   string `json:"symbol"`
}

package token

import (
	"math/big"

	"github.com/loomnetwork/go-loom/auth"
)

type Erc20 interface {
	BalanceOf(hexAddr string) (*big.Int, error)
	Name() string
	Decimals() uint8
	Symbol() string
	Approve(signer auth.Signer, to string, amount *big.Int) (string, error)
	Transfer(signer auth.Signer, to string, amount *big.Int) (string, error)
}

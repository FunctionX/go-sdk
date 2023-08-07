package loadtest

import (
	"github.com/functionx/go-sdk/cosmos/auth"
	"github.com/functionx/go-sdk/cosmos/types"
)

type RPCClient interface {
	GetAddressPrefix() (string, error)
	QueryAccount(address string) (auth.AccountI, error)
	GetChainId() (string, error)
	QuerySupply() (types.Coins, error)
}

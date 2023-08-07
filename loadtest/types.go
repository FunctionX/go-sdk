package loadtest

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/pkg/errors"

	"github.com/functionx/go-sdk/cosmos"
	"github.com/functionx/go-sdk/cosmos/client"
	"github.com/functionx/go-sdk/cosmos/grpc"
	cosmostypes "github.com/functionx/go-sdk/cosmos/types"
)

type BaseInfo struct {
	Accounts *Accounts
	ChainID  string
	GasPrice cosmostypes.DecCoin
	GasLimit int64
	Memo     string
}

func newBaseInfo(accounts *Accounts, chainId, denom string) *BaseInfo {
	return &BaseInfo{
		Accounts: accounts,
		ChainID:  chainId,
		GasPrice: cosmostypes.NewDecCoinFromDec(denom, sdkmath.LegacyNewDec(0)),
		GasLimit: 100_000,
		Memo:     "",
	}
}

func (i *BaseInfo) GetDenom() string {
	return i.GasPrice.Denom
}

func NewBaseInfo(str string, keyDir string) (*BaseInfo, error) {
	if strings.HasSuffix(str, "config/genesis.json") {
		return NewBaseInfoFromGenesis(str, keyDir)
	}
	if strings.Contains(str, "://") {
		newClient := client.NewClient(str, cosmos.NewProtoCodec())
		grpcCli := grpc.NewClient(context.Background(), newClient)
		return NewBaseInfoFromClient(grpcCli, keyDir)
	} else {
		return nil, errors.New("invalid base info")
	}
}

func NewBaseInfoFromClient(client RPCClient, keyDir string) (*BaseInfo, error) {
	accounts, err := NewAccounts(client, keyDir)
	if err != nil {
		return nil, err
	}
	chainID, err := client.GetChainId()
	if err != nil {
		return nil, err
	}
	supply, err := client.QuerySupply()
	if err != nil {
		return nil, err
	}
	denom := supply[0].Denom
	return newBaseInfo(accounts, chainID, denom), nil
}

func NewBaseInfoFromGenesis(genesisPath string, keyDir string) (*BaseInfo, error) {
	accounts, err := NewAccountFromGenesis(genesisPath, keyDir)
	if err != nil {
		return nil, err
	}
	genesisFile, err := os.ReadFile(genesisPath)
	if err != nil {
		return nil, err
	}
	var genesis struct {
		ChainId  string `json:"chain_id"`
		AppState struct {
			Mint struct {
				Params struct {
					MintDenom string `json:"mint_denom"`
				} `json:"params"`
			} `json:"mint"`
			Staking struct {
				Params struct {
					BondDenom string `json:"bond_denom"`
				} `json:"params"`
			} `json:"staking"`
		} `json:"app_state"`
	}
	if err = json.Unmarshal(genesisFile, &genesis); err != nil {
		return nil, err
	}
	denom := genesis.AppState.Staking.Params.BondDenom
	return newBaseInfo(accounts, genesis.ChainId, denom), nil
}

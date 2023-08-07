package grpc

import (
	"context"

	"github.com/functionx/go-sdk/cosmos/codec"
	cosmostypes "github.com/functionx/go-sdk/cosmos/types"
	"github.com/functionx/go-sdk/osmosis"
	"github.com/functionx/go-sdk/osmosis/gamm"
	"github.com/functionx/go-sdk/osmosis/poolmanager"
)

type OsmosisClient struct {
	*Client
	poolCodec *codec.ProtoCodec
}

func NewOsmosisClient(ctx context.Context, config Config) (*OsmosisClient, error) {
	cosmosClient, err := newCosmosClient(ctx, config)
	if err != nil {
		return nil, err
	}
	return &OsmosisClient{Client: cosmosClient, poolCodec: osmosis.NewProtoCodec()}, nil
}

func (cli *OsmosisClient) PoolManagerQuery() poolmanager.QueryClient {
	return poolmanager.NewQueryClient(cli.ClientConn)
}

func (cli *OsmosisClient) GammQuery() gamm.QueryClient {
	return gamm.NewQueryClient(cli.ClientConn)
}

func (cli *OsmosisClient) CalcExitPoolCoinsFromShares(ctx context.Context, poolId uint64, shareAmount cosmostypes.Int) (cosmostypes.Coins, error) {
	response, err := cli.GammQuery().CalcExitPoolCoinsFromShares(ctx, &gamm.QueryCalcExitPoolCoinsFromSharesRequest{PoolId: poolId, ShareInAmount: shareAmount})
	if err != nil {
		return nil, err
	}
	return response.GetTokensOut(), nil
}

func (cli *OsmosisClient) EstimateSwapExactAmountIn(ctx context.Context, tokenIn string, router []poolmanager.SwapAmountInRoute) (cosmostypes.Int, error) {
	response, err := cli.PoolManagerQuery().EstimateSwapExactAmountIn(ctx, &poolmanager.EstimateSwapExactAmountInRequest{
		TokenIn: tokenIn,
		Routes:  router,
	})
	if err != nil {
		return cosmostypes.Int{}, err
	}
	return response.TokenOutAmount, nil
}

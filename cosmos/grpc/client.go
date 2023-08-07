package grpc

import (
	"context"

	"github.com/pkg/errors"

	"github.com/functionx/go-sdk/cosmos/auth"
	"github.com/functionx/go-sdk/cosmos/client/tmservice"
	"github.com/functionx/go-sdk/cosmos/types"
	"github.com/functionx/go-sdk/cosmos/types/tx"
	"github.com/functionx/go-sdk/log"
)

type RPCClient interface {
	GetChainId() (string, error)
	QueryAccount(address string) (auth.AccountI, error)
	QueryBalance(address string, denom string) (types.Coin, error)
	GetGasPrices() (types.DecCoins, error)
	EstimatingGas(raw *tx.TxRaw) (*types.GasInfo, error)
	BroadcastTx(txRaw *tx.TxRaw, mode ...tx.BroadcastMode) (*types.TxResponse, error)
	TxByHash(txHash string) (*types.TxResponse, error)
	TxSearch(event []string) ([]*types.TxResponse, error)
	GetLatestBlock() (*tmservice.Block, error)

	Close() error
}

type client struct {
	*Client
	logger log.Logger
	config Config
}

func NewRPCClient(ctx context.Context, logger log.Logger, config Config) (RPCClient, error) {
	cosmosClient, err := newCosmosClient(ctx, config)
	if err != nil {
		return nil, err
	}
	rpcClient := &client{
		Client: cosmosClient,
		logger: logger.With("module", "cosmos-grpc"),
		config: config,
	}
	return rpcClient, rpcClient.check(ctx)
}

func newCosmosClient(ctx context.Context, config Config) (*Client, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()
	conn, err := NewGrpcConn(timeoutCtx, config.RpcUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "dial grpc url: %s", config.RpcUrl)
	}
	return NewClient(ctx, conn), nil
}

func (c *client) check(_ context.Context) error {
	id, err := c.GetChainId()
	if err != nil {
		return errors.Wrap(err, "get chain id")
	}
	if id != c.chainId {
		return errors.Errorf("chain id mismatch, expected %s, got %s", id, c.chainId)
	}

	prefix, err := c.GetAddressPrefix()
	if err != nil {
		return errors.Wrap(err, "get address prefix")
	}
	if prefix != c.config.AddressPrefix {
		return errors.Errorf("address prefix mismatch, expected %s, got %s", prefix, c.config.AddressPrefix)
	}
	return nil
}

func (c *client) Close() error {
	if err := c.Client.Close(); err != nil {
		return errors.Wrap(err, "close cosmos grpc client")
	}
	return nil
}

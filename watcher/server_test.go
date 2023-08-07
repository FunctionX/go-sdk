package watcher

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"

	"github.com/functionx/go-sdk/cosmos"
	"github.com/functionx/go-sdk/cosmos/client/tmservice"
	"github.com/functionx/go-sdk/cosmos/types"
	"github.com/functionx/go-sdk/log"
)

func TestNewServer(t *testing.T) {
	cfg := NewDefConfig()
	codec := cosmos.NewProtoCodec()

	logger, err := log.NewLogger(log.FormatConsole, "info")
	assert.NoError(t, err)
	server := NewServer(logger, cfg, codec)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)

	server.client = MockRpcClient{}
	assert.NoError(t, server.Start(group, ctx))

	<-ctx.Done()
	assert.NoError(t, server.Close())
	assert.Error(t, group.Wait())
}

var _ RPCClient = (*MockRpcClient)(nil)

type MockRpcClient struct{}

func (m MockRpcClient) TxByHash(txHash string) (*types.TxResponse, error) {
	return nil, assert.AnError
}

func (m MockRpcClient) GetLatestBlock() (*tmservice.Block, error) {
	return nil, assert.AnError
}

func (m MockRpcClient) Close() error {
	return nil
}

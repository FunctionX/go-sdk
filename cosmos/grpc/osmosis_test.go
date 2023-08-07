package grpc_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/functionx/go-sdk/cosmos/grpc"
)

func TestOsmosisClient(t *testing.T) {
	if os.Getenv("GRPC_CLIENT_TESTS") != "1" {
		t.Skip("skipping grpc client test")
	}

	osmoCli, err := grpc.NewOsmosisClient(context.Background(), grpc.Config{
		ChainId:       "osmo-test-5",
		RpcUrl:        "https://grpc.osmotest5.osmosis.zone",
		Timeout:       time.Second * 10,
		AddressPrefix: "osmo",
	})
	assert.NoError(t, err)
	chainID, err := osmoCli.GetChainId()
	assert.NoError(t, err)
	assert.Equal(t, "osmo-test-5", chainID)
}

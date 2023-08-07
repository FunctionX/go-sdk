package client

import (
	"context"
	"testing"

	"github.com/functionx/go-sdk/cosmos"
	"github.com/functionx/go-sdk/cosmos/grpc"
	"github.com/functionx/go-sdk/tendermint/abci"
)

func Test_NewClient(t *testing.T) {
	t.Skip("skip client test")

	newClient := NewClient("http://127.0.0.1:26657", cosmos.NewProtoCodec())
	queryABCI, err := newClient.QueryABCI(context.Background(), abci.RequestQuery{
		Path: "/app/version",
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	t.Log(string(queryABCI.Value))

	grpcClient := grpc.NewClient(context.Background(), newClient)
	chainId, err := grpcClient.GetChainId()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	t.Log(chainId)

	prefix, err := grpcClient.GetAddressPrefix()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	t.Log(prefix)

	account, err := grpcClient.QueryAccount("0x433d0d9f2900229069F44Ff0F7ea4CE40d54F639")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	t.Log(account.String())
}

package loadtest

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/informalsystems/tm-load-test/pkg/loadtest"

	"github.com/functionx/go-sdk/cosmos"
	csomosclient "github.com/functionx/go-sdk/cosmos/client"
	"github.com/functionx/go-sdk/cosmos/grpc"
	"github.com/functionx/go-sdk/cosmos/types/tx"
)

func Test_SendMsgOrder(t *testing.T) {
	t.Skip("skip load test")
	// grpcClient, err := grpc.DialContext(context.Background(), "http://127.0.0.1:9090")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	newClient := csomosclient.NewClient("http://127.0.0.1:26657", cosmos.NewProtoCodec())
	grpcClient := grpc.NewClient(context.Background(), newClient)
	keyDir := filepath.Join(os.ExpandEnv("$HOME"), "test_accounts")
	baseInfo, err := NewBaseInfoFromClient(grpcClient, keyDir)
	if err != nil {
		t.Fatal(err)
	}
	msgOrderClientFactory := NewMsgOrderClientFactory(baseInfo)
	client, err := msgOrderClientFactory.NewClient(loadtest.Config{})
	if err != nil {
		t.Fatal(err)
	}
	rawTx, err := client.GenerateTx()
	if err != nil {
		t.Fatal(err)
	}
	txResp, err := grpcClient.BroadcastTxBytes(rawTx, tx.BroadcastMode_BROADCAST_MODE_SYNC)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txResp.Code, txResp.String())
}

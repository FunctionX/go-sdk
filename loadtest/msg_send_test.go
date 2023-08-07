package loadtest

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/gogo/protobuf/proto"
	"github.com/informalsystems/tm-load-test/pkg/loadtest"
	"github.com/stretchr/testify/assert"

	"github.com/functionx/go-sdk/cosmos"
	csomosclient "github.com/functionx/go-sdk/cosmos/client"
	"github.com/functionx/go-sdk/cosmos/grpc"
	cosmostypes "github.com/functionx/go-sdk/cosmos/types"
	typestx "github.com/functionx/go-sdk/cosmos/types/tx"
)

func Test_MsgSend(t *testing.T) {
	t.Skip("skip load test")

	newClient := csomosclient.NewClient("http://127.0.0.1:26657", cosmos.NewProtoCodec())
	grpcClient := grpc.NewClient(context.Background(), newClient)
	keyDir := filepath.Join(os.ExpandEnv("$HOME"), "sei_test_accounts")
	genesisFilePath := filepath.Join(os.ExpandEnv("$HOME"), ".sei", "config", "genesis.json")
	baseInfo, err := NewBaseInfoFromGenesis(genesisFilePath, keyDir)
	if err != nil {
		t.Fatal(err)
	}
	msgSendClientFactory := NewMsgSendClientFactory(baseInfo, baseInfo.GetDenom())
	t.Logf(msgSendClientFactory.GasPrice.String())
	client, err := msgSendClientFactory.NewClient(loadtest.Config{})
	if err != nil {
		t.Fatal(err)
	}
	rawTx, err := client.GenerateTx()
	if err != nil {
		t.Fatal(err)
	}
	txResp, err := grpcClient.BroadcastTxBytes(rawTx, typestx.BroadcastMode_BROADCAST_MODE_SYNC)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txResp.Code, txResp.String())
}

// go test -bench ^Benchmark_NewMsgSendTxsAndMarshal$ -benchtime 10s -count 1 -cpu 1 -run=^$ -benchmem
func Benchmark_NewMsgSendTxsAndMarshal(b *testing.B) {
	b.Skip("skip load test")

	homeDir := os.ExpandEnv("$HOME")
	keyDir := filepath.Join(homeDir, "test_accounts")

	genesisFilePath := filepath.Join(homeDir, ".simapp", "config", "genesis.json")
	accounts, err := NewAccountFromGenesis(genesisFilePath, keyDir)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = createTestTx(accounts); err != nil {
			b.Fatal(err)
		}
	}
}

// go test -bench ^Benchmark_NewMsgSendTxsAndUnmarshal$ -benchtime 10s -count 1 -cpu 1 -run=^$ -benchmem
func Benchmark_NewMsgSendTxsAndUnmarshal(b *testing.B) {
	b.Skip("skip load test")

	homeDir := os.ExpandEnv("$HOME")
	keyDir := filepath.Join(homeDir, "test_accounts")

	genesisFilePath := filepath.Join(homeDir, ".simapp", "config", "genesis.json")
	accounts, err := NewAccountFromGenesis(genesisFilePath, keyDir)
	if err != nil {
		b.Fatal(err)
	}
	b.Log(accounts.Len())
	txs, err := createTestTx(accounts)
	if err != nil {
		b.Fatal(err)
	}
	codec := cosmos.NewProtoCodec()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(txs); j++ {
			var txRawBytes []byte
			txRawBytes, err = base64.StdEncoding.DecodeString(txs[j])
			if err != nil {
				b.Fatal(err)
			}
			var txRaw typestx.TxRaw
			if err = codec.Unmarshal(txRawBytes, &txRaw); err != nil {
				b.Fatal(err)
			}

			var body typestx.TxBody
			if err = codec.Unmarshal(txRaw.BodyBytes, &body); err != nil {
				b.Fatal(err)
			}

			var authInfo typestx.AuthInfo
			if err = codec.Unmarshal(txRaw.AuthInfoBytes, &authInfo); err != nil {
				b.Fatal(err)
			}
			for n, info := range authInfo.SignerInfos {
				var signatures []byte
				signatures, err = proto.Marshal(&typestx.SignDoc{
					BodyBytes:     txRaw.BodyBytes,
					AuthInfoBytes: txRaw.AuthInfoBytes,
					ChainId:       "cosmos",
					AccountNumber: accounts.accounts[j].AccountNumber,
				})
				if err != nil {
					b.Fatal(err)
				}
				if !info.GetPubKey().VerifySignature(signatures, txRaw.GetSignatures()[n]) {
					panic("verify signature failed")
				}
			}
		}
	}
}

func Test_MsgSend2(t *testing.T) {
	t.Skip("skip load test")

	homeDir := os.ExpandEnv("$HOME")
	keyDir := filepath.Join(homeDir, "test_accounts")

	err := CreateGenesisAccounts("sei", 10, keyDir)
	assert.NoError(t, err)

	genesisFilePath := filepath.Join(homeDir, ".simapp", "config", "genesis.json")
	accounts, err := NewAccountFromGenesis(genesisFilePath, keyDir)
	if err != nil {
		t.Fatal(err)
	}
	txs, err := createTestTx(accounts)
	if err != nil {
		t.Fatal(err)
	}
	if err = writeTxsToFile(txs); err != nil {
		t.Fatal(err)
	}
}

func createTestTx(accounts *Accounts) ([]string, error) {
	factory := NewMsgSendClientFactory(&BaseInfo{
		Accounts: accounts,
		ChainID:  "cosmos",
		// GasPrice: cosmostypes.NewDecCoinFromDec("stake", sdkmath.LegacyNewDecWithPrec(1, 2)),
		GasPrice: cosmostypes.NewDecCoinFromDec("stake", sdkmath.LegacyNewDec(0)),
		GasLimit: 120000,
	}, "stake")

	txCount := accounts.Len()
	txs := make([]string, txCount)
	for i := 0; i < txCount; i++ {
		txRaw, err := factory.GenerateTx()
		if err != nil {
			return nil, err
		}
		txs[i] = base64.StdEncoding.EncodeToString(txRaw)
	}

	return txs, nil
}

func writeTxsToFile(txs []string) error {
	txData, err := json.Marshal(txs)
	if err != nil {
		return err
	}
	if err = os.RemoveAll(filepath.Join(os.ExpandEnv("$HOME"), "txs.json")); err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(os.ExpandEnv("$HOME"), "txs.json"), txData, 0o600); err != nil {
		return err
	}
	return nil
}

package grpc_test

import (
	"context"
	"os"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/functionx/go-sdk/cosmos/auth"
	"github.com/functionx/go-sdk/cosmos/bank"
	ethsecp256k12 "github.com/functionx/go-sdk/cosmos/crypto/ethsecp256k1"
	"github.com/functionx/go-sdk/cosmos/grpc"
	"github.com/functionx/go-sdk/cosmos/types"
	"github.com/functionx/go-sdk/cosmos/types/tx"
)

const TestAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

func TestClient(t *testing.T) {
	if os.Getenv("GRPC_CLIENT_TESTS") != "1" {
		t.Skip("skipping grpc client test")
	}

	const expectedChainId = "fxcore"
	client, err := grpc.DialContext(context.Background(), "http://localhost:9090")
	assert.NoError(t, err)

	chainId, err := client.GetChainId()
	assert.NoError(t, err)
	assert.Equal(t, expectedChainId, chainId)

	height, err := client.GetBlockHeight()
	assert.NoError(t, err)
	assert.True(t, height > 0)

	version, err := client.AppVersion()
	assert.NoError(t, err)
	assert.NotEmpty(t, version)

	supply, err := client.QuerySupply()
	assert.NoError(t, err)
	assert.True(t, supply.IsAllPositive())

	prefix, err := client.GetAddressPrefix()
	assert.NoError(t, err)
	assert.NotEmpty(t, prefix)

	prices, err := client.GetGasPrices()
	assert.NoError(t, err)
	assert.True(t, prices.IsAllPositive())

	syncing, err := client.GetSyncing()
	assert.NoError(t, err)
	assert.False(t, syncing)

	hexAddress := common.HexToAddress(TestAddress)
	accAddress := types.NewAccAddress(hexAddress.Bytes(), prefix)

	account, err := client.QueryAccount(accAddress.String())
	assert.NoError(t, err)

	assert.NotNil(t, account.GetPubKey())
	pubKey, err := types.NewAnyWithValue(&ethsecp256k12.PubKey{Key: account.GetPubKey().Bytes()})
	assert.NoError(t, err)

	assert.Equal(t, &auth.BaseAccount{
		Address:       accAddress.String(),
		PubKey:        pubKey,
		AccountNumber: 0,
		Sequence:      1,
	}, account)

	balances, err := client.QueryBalances(accAddress.String())
	assert.NoError(t, err)
	assert.True(t, balances.IsAllPositive())
}

func TestClient_Tx(t *testing.T) {
	if os.Getenv("GRPC_CLIENT_TESTS") != "1" {
		t.Skip("skipping grpc client test")
	}

	client, err := grpc.DialContext(context.Background(), "http://localhost:9090")
	assert.NoError(t, err)

	keyjson := `{"address":"f39fd6e51aad88f6f4ce6ab8827279cfffb92266","crypto":{"cipher":"aes-128-ctr","ciphertext":"d4164b3b560085ecec3063b3a6b8c15202775844be8b596cfb45d6025dd06a39","cipherparams":{"iv":"1fe40a9b659a97e6ec45b8bf1b7fe2cb"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"6247060bd8c08b14b8aa9877195010e915a28da098694e0a6deedf51160b8ba2"},"mac":"240038f0ed8235bac90b8472dfae2ae562836ecfd25a4c71b9ccc4009b0ad3ad"},"id":"00000000-0000-0000-0000-000000000000","version":3}`
	key, err := keystore.DecryptKey([]byte(keyjson), "12345678")
	assert.NoError(t, err)
	privKey := ethsecp256k12.NewPrivKey(key.PrivateKey)

	txRaw, err := client.BuildTx(
		privKey,
		[]types.Msg{
			&bank.MsgSend{
				FromAddress: types.NewAccAddress(privKey.PubKey().Address().Bytes(), "fx").String(),
				ToAddress:   types.NewAccAddress(privKey.PubKey().Address().Bytes(), "fx").String(),
				Amount:      types.NewCoins(types.NewCoin("FX", sdkmath.NewInt(1))),
			},
		},
	)
	assert.NoError(t, err)

	gas, err := client.EstimatingGas(txRaw)
	assert.NoError(t, err)
	assert.True(t, gas.GasUsed < 110000)

	txResponse, err := client.BroadcastTx(txRaw)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), txResponse.Code)
	assert.True(t, txResponse.GasUsed < 110000)

	time.Sleep(500 * time.Millisecond)
	txRes, err := client.TxByHash(txResponse.TxHash)
	assert.NoError(t, err)
	txRes.Tx = nil
	txRes.Timestamp = ""
	assert.Equal(t, txResponse, txRes)
}

func TestGetTxsEvent(t *testing.T) {
	if os.Getenv("GRPC_CLIENT_TESTS") != "1" {
		t.Skip("skipping grpc client test")
	}
	client, err := grpc.DialContext(context.Background(), "https://fx-grpc.functionx.io:9090")
	assert.NoError(t, err)
	//	eventRequest := &tx.GetTxsEventRequest{Events: []string{"send_packet.packet_src_channel='channel-11'"}, OrderBy: tx.OrderBy_ORDER_BY_DESC}
	eventRequest := &tx.GetTxsEventRequest{Events: []string{"recv_packet.packet_src_channel='channel-784'", "recv_packet.packet_sequence='3'"}, OrderBy: tx.OrderBy_ORDER_BY_DESC}
	//	eventRequest := &tx.GetTxsEventRequest{Events: []string{"ethereum_tx.ethereumTxHash='0x77cb62c98caf13cf347a47bee07ae0a117ea854f85d1d03603265ea3bc7cf627'"}, OrderBy: tx.OrderBy_ORDER_BY_DESC}
	txsEvent, err := client.ServiceClient().GetTxsEvent(context.Background(), eventRequest)
	assert.NoError(t, err)
	txResponses := txsEvent.TxResponses
	for _, response := range txResponses {
		t.Logf("%v", response.TxHash)
	}
}

package grpc

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	gogogrpc "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/google"
	"google.golang.org/grpc/metadata"

	"github.com/functionx/go-sdk/bech32"
	"github.com/functionx/go-sdk/cosmos"
	"github.com/functionx/go-sdk/cosmos/auth"
	"github.com/functionx/go-sdk/cosmos/bank"
	"github.com/functionx/go-sdk/cosmos/client/node"
	"github.com/functionx/go-sdk/cosmos/client/tmservice"
	"github.com/functionx/go-sdk/cosmos/codec"
	cosmoscrypto "github.com/functionx/go-sdk/cosmos/crypto"
	"github.com/functionx/go-sdk/cosmos/types"
	"github.com/functionx/go-sdk/cosmos/types/tx"
)

const (
	// BlockHeightHeader is the gRPC header for block height.
	BlockHeightHeader = "x-cosmos-block-height"
)

type ClientConn interface {
	gogogrpc.ClientConn
	Close() error
}

type Client struct {
	chainId      string
	addrPrefix   string
	gasPrices    types.DecCoins
	ctx          context.Context
	accountCodec *codec.ProtoCodec
	ClientConn
}

func NewGrpcConn(ctx context.Context, rawUrl string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	_url := u.Host
	if u.Port() == "" {
		if u.Scheme == "https" {
			_url = u.Host + ":443"
		} else {
			_url = u.Host + ":80"
		}
	}
	if u.Scheme == "https" {
		opts = append(opts, grpc.WithCredentialsBundle(google.NewDefaultCredentials()))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	return grpc.DialContext(ctx, _url, opts...)
}

func NewClient(ctx context.Context, conn ClientConn) *Client {
	return &Client{ClientConn: conn, ctx: ctx, accountCodec: cosmos.NewProtoCodec()}
}

func DialContext(ctx context.Context, rawUrl string) (*Client, error) {
	grpcConn, err := NewGrpcConn(ctx, rawUrl)
	if err != nil {
		return nil, err
	}
	return NewClient(ctx, grpcConn), nil
}

func (cli *Client) WithContext(ctx context.Context) *Client {
	return &Client{
		ClientConn:   cli.ClientConn,
		chainId:      cli.chainId,
		gasPrices:    cli.gasPrices,
		addrPrefix:   cli.addrPrefix,
		ctx:          ctx,
		accountCodec: cli.accountCodec,
	}
}

func (cli *Client) WithBlockHeight(height int64) *Client {
	ctx := metadata.AppendToOutgoingContext(cli.ctx, BlockHeightHeader, strconv.FormatInt(height, 10))
	return &Client{
		ClientConn:   cli.ClientConn,
		chainId:      cli.chainId,
		gasPrices:    cli.gasPrices,
		addrPrefix:   cli.addrPrefix,
		ctx:          ctx,
		accountCodec: cli.accountCodec,
	}
}

func (cli *Client) AuthQuery() auth.QueryClient {
	return auth.NewQueryClient(cli.ClientConn)
}

func (cli *Client) BankQuery() bank.QueryClient {
	return bank.NewQueryClient(cli.ClientConn)
}

func (cli *Client) ServiceClient() tx.ServiceClient {
	return tx.NewServiceClient(cli.ClientConn)
}

func (cli *Client) TMServiceClient() tmservice.ServiceClient {
	return tmservice.NewServiceClient(cli.ClientConn)
}

func (cli *Client) AppVersion() (string, error) {
	info, err := cli.TMServiceClient().GetNodeInfo(cli.ctx, &tmservice.GetNodeInfoRequest{})
	if err != nil {
		return "", err
	}
	return info.GetApplicationVersion().GetVersion(), nil
}

func (cli *Client) QueryAccount(address string) (auth.AccountI, error) {
	response, err := cli.AuthQuery().Account(cli.ctx, &auth.QueryAccountRequest{
		Address: address,
	})
	if err != nil {
		return nil, err
	}
	var account auth.AccountI
	if err = cli.accountCodec.UnpackAny(response.GetAccount(), &account); err != nil {
		return nil, err
	}
	return account, nil
}

func (cli *Client) QueryBalance(address string, denom string) (types.Coin, error) {
	response, err := cli.BankQuery().Balance(cli.ctx, &bank.QueryBalanceRequest{
		Address: address,
		Denom:   denom,
	})
	if err != nil {
		return types.Coin{}, err
	}
	return *response.GetBalance(), nil
}

func (cli *Client) QueryBalances(address string) (types.Coins, error) {
	response, err := cli.BankQuery().AllBalances(cli.ctx, &bank.QueryAllBalancesRequest{
		Address: address,
	})
	if err != nil {
		return nil, err
	}
	return response.GetBalances(), nil
}

func (cli *Client) QuerySupply() (types.Coins, error) {
	response, err := cli.BankQuery().TotalSupply(cli.ctx, &bank.QueryTotalSupplyRequest{})
	if err != nil {
		return nil, err
	}
	return response.GetSupply(), nil
}

func (cli *Client) GetBlockHeight() (int64, error) {
	response, err := cli.TMServiceClient().GetLatestBlock(cli.ctx, &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return 0, err
	}
	return response.GetSdkBlock().GetHeader().Height, nil
}

func (cli *Client) GetChainId() (string, error) {
	if len(cli.chainId) > 0 {
		return cli.chainId, nil
	}
	response, err := cli.TMServiceClient().GetLatestBlock(cli.ctx, &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return "", err
	}
	var chainId string
	if response.GetSdkBlock() != nil {
		chainId = response.GetSdkBlock().GetHeader().ChainID
	} else {
		chainId = response.GetBlock().GetHeader().ChainID
	}
	cli.chainId = chainId
	return chainId, nil
}

func (cli *Client) GetBlockTimeInterval() (time.Duration, error) {
	tmClient := cli.TMServiceClient()
	response1, err := tmClient.GetLatestBlock(cli.ctx, &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return 0, err
	}
	if response1.GetSdkBlock().GetHeader().Height <= 1 {
		return 0, fmt.Errorf("please try again later, the current block height is less than 1")
	}
	response2, err := tmClient.GetBlockByHeight(cli.ctx, &tmservice.GetBlockByHeightRequest{
		Height: response1.GetSdkBlock().GetHeader().Height - 1,
	})
	if err != nil {
		return 0, err
	}
	return response1.GetSdkBlock().GetHeader().Time.Sub(response2.GetSdkBlock().GetHeader().Time), nil
}

func (cli *Client) GetLatestBlock() (*tmservice.Block, error) {
	response, err := cli.TMServiceClient().GetLatestBlock(cli.ctx, &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return nil, err
	}
	return response.GetSdkBlock(), nil
}

func (cli *Client) GetBlockByHeight(blockHeight int64) (*tmservice.Block, error) {
	response, err := cli.TMServiceClient().GetBlockByHeight(cli.ctx, &tmservice.GetBlockByHeightRequest{
		Height: blockHeight,
	})
	if err != nil {
		return nil, err
	}
	return response.GetSdkBlock(), nil
}

func (cli *Client) GetStatusByTx(txHash string) (*tx.GetTxResponse, error) {
	response, err := cli.ServiceClient().GetTx(cli.ctx, &tx.GetTxRequest{
		Hash: txHash,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (cli *Client) GetGasPrices() (types.DecCoins, error) {
	if len(cli.gasPrices) > 0 {
		return cli.gasPrices, nil
	}
	response, err := node.NewServiceClient(cli).Config(cli.ctx, &node.ConfigRequest{})
	if err != nil {
		return nil, err
	}
	coins, err := types.ParseDecCoins(response.GetMinimumGasPrice())
	if err != nil {
		return nil, err
	}
	cli.gasPrices = coins
	return coins, nil
}

func (cli *Client) GetAddressPrefix() (string, error) {
	if len(cli.addrPrefix) > 0 {
		return cli.addrPrefix, nil
	}
	response, err := cli.TMServiceClient().GetLatestValidatorSet(cli.ctx, &tmservice.GetLatestValidatorSetRequest{})
	if err != nil {
		return "", err
	}
	if len(response.GetValidators()) == 0 {
		return "", fmt.Errorf("no found validator")
	}
	address := response.GetValidators()[0].GetAddress()
	if strings.HasPrefix(address, "0x") {
		return "0x", nil
	}
	prefix, _, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return "", err
	}
	valConPrefix := "valcons"
	if strings.HasSuffix(prefix, valConPrefix) {
		cli.addrPrefix = prefix[:len(prefix)-len(valConPrefix)]
		return cli.addrPrefix, nil
	}
	return "", fmt.Errorf("no found address prefix")
}

func (cli *Client) GetSyncing() (bool, error) {
	response, err := cli.TMServiceClient().GetSyncing(cli.ctx, &tmservice.GetSyncingRequest{})
	if err != nil {
		return false, err
	}
	return response.Syncing, nil
}

func (cli *Client) GetNodeInfo() (*tmservice.VersionInfo, error) {
	response, err := cli.TMServiceClient().GetNodeInfo(cli.ctx, &tmservice.GetNodeInfoRequest{})
	if err != nil {
		return nil, err
	}
	return response.GetApplicationVersion(), nil
}

func (cli *Client) GetConsensusValidators() ([]*tmservice.Validator, error) {
	response, err := cli.TMServiceClient().GetLatestValidatorSet(cli.ctx, &tmservice.GetLatestValidatorSetRequest{})
	if err != nil {
		return nil, err
	}
	return response.GetValidators(), nil
}

func (cli *Client) EstimatingGas(raw *tx.TxRaw) (*types.GasInfo, error) {
	txBytes, err := proto.Marshal(raw)
	if err != nil {
		return nil, err
	}
	response, err := cli.ServiceClient().Simulate(cli.ctx, &tx.SimulateRequest{TxBytes: txBytes})
	if err != nil {
		return nil, err
	}
	return response.GetGasInfo(), nil
}

func (cli *Client) BuildTx(privKey cosmoscrypto.PrivKey, msgs []types.Msg) (*tx.TxRaw, error) {
	return cosmos.BuildTx(cli, privKey, msgs)
}

func (cli *Client) BroadcastTxOk(txRaw *tx.TxRaw, mode ...tx.BroadcastMode) (*types.TxResponse, error) {
	broadcastTx, err := cli.BroadcastTx(txRaw, mode...)
	if err != nil {
		return nil, err
	}
	if broadcastTx.Code != 0 {
		return nil, fmt.Errorf(broadcastTx.RawLog)
	}
	return broadcastTx, nil
}

func (cli *Client) BroadcastTx(txRaw *tx.TxRaw, mode ...tx.BroadcastMode) (*types.TxResponse, error) {
	txBytes, err := proto.Marshal(txRaw)
	if err != nil {
		return nil, err
	}
	defaultMode := tx.BroadcastMode_BROADCAST_MODE_BLOCK
	if len(mode) > 0 {
		defaultMode = mode[0]
	}

	_, err = proto.Marshal(&tx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    defaultMode,
	})
	if err != nil {
		return nil, err
	}
	broadcastTxResponse, err := cli.ServiceClient().BroadcastTx(cli.ctx, &tx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    defaultMode,
	})
	if err != nil {
		return nil, err
	}
	return broadcastTxResponse.GetTxResponse(), nil
}

func (cli *Client) BroadcastTxBytes(txBytes []byte, mode ...tx.BroadcastMode) (*types.TxResponse, error) {
	defaultMode := tx.BroadcastMode_BROADCAST_MODE_BLOCK
	if len(mode) > 0 {
		defaultMode = mode[0]
	}
	_, err := proto.Marshal(&tx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    defaultMode,
	})
	if err != nil {
		return nil, err
	}
	broadcastTxResponse, err := cli.ServiceClient().BroadcastTx(cli.ctx, &tx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    defaultMode,
	})
	if err != nil {
		return nil, err
	}
	return broadcastTxResponse.GetTxResponse(), nil
}

func (cli *Client) TxByHash(txHash string) (*types.TxResponse, error) {
	resp, err := cli.ServiceClient().GetTx(cli.ctx, &tx.GetTxRequest{Hash: txHash})
	if err != nil {
		return nil, err
	}
	return resp.GetTxResponse(), nil
}

func (cli *Client) TxSearch(events []string) ([]*types.TxResponse, error) {
	resp, err := cli.ServiceClient().GetTxsEvent(cli.ctx, &tx.GetTxsEventRequest{Events: events, OrderBy: tx.OrderBy_ORDER_BY_DESC})
	if err != nil {
		return nil, err
	}
	return resp.TxResponses, nil
}

func (cli *Client) BuildTxV1(privKey cosmoscrypto.PrivKey, msgs []types.Msg, gasLimit int64, memo string, timeout uint64) (*tx.TxRaw, error) {
	prefix, err := cli.GetAddressPrefix()
	if err != nil {
		return nil, err
	}
	from, err := bech32.ConvertAndEncode(prefix, privKey.PubKey().Address())
	if err != nil {
		return nil, err
	}
	account, err := cli.QueryAccount(from)
	if err != nil {
		return nil, err
	}
	chainId, err := cli.GetChainId()
	if err != nil {
		return nil, err
	}
	var gasPrice types.DecCoin
	gasPrices, err := cli.GetGasPrices()
	if err != nil {
		return nil, err
	}
	if len(gasPrices) > 0 {
		gasPrice = gasPrices[0]
	}

	txRaw, err := cosmos.BuildTxV1(chainId, account.GetSequence(), account.GetAccountNumber(), privKey, msgs, gasPrice, gasLimit, memo, timeout)
	if err != nil {
		return nil, err
	}
	estimatingGas, err := cli.EstimatingGas(txRaw)
	if err != nil {
		return nil, err
	}
	if estimatingGas.GetGasUsed() > uint64(gasLimit) {
		gasLimit = int64(estimatingGas.GetGasUsed()) + (int64(estimatingGas.GetGasUsed()) * 2 / 10)
	}
	return cosmos.BuildTxV1(chainId, account.GetSequence(), account.GetAccountNumber(), privKey, msgs, gasPrice, gasLimit, memo, timeout)
}

package client

import (
	"context"
	"errors"

	"github.com/functionx/go-sdk/cosmos/codec"
	cosmostypes "github.com/functionx/go-sdk/cosmos/types"
	"github.com/functionx/go-sdk/cosmos/types/tx"
	"github.com/functionx/go-sdk/tendermint/jsonrpc"
)

type Client struct {
	jsonRPC *jsonrpc.JsonRPC
	height  int64
	cdc     codec.ProtoCodecMarshaler
}

func NewClient(url string, cdc codec.ProtoCodecMarshaler) *Client {
	return &Client{
		jsonRPC: jsonrpc.NewJsonRPC(url),
		cdc:     cdc,
	}
}

// TxServiceBroadcast is a helper function to broadcast a Tx with the correct gRPC types
// from the tx service. Calls `clientCtx.BroadcastTx` under the hood.
func (c *Client) TxServiceBroadcast(ctx context.Context, req *tx.BroadcastTxRequest) (*tx.BroadcastTxResponse, error) {
	if req == nil || req.TxBytes == nil {
		return nil, errors.New("invalid empty tx")
	}

	var method string
	switch req.Mode {
	case tx.BroadcastMode_BROADCAST_MODE_BLOCK:
		method = "broadcast_tx_commit"
	case tx.BroadcastMode_BROADCAST_MODE_SYNC:
		method = "broadcast_tx_sync"
	case tx.BroadcastMode_BROADCAST_MODE_ASYNC:
		method = "broadcast_tx_async"
	default:
		return nil, errors.New("invalid broadcast mode")
	}

	var txResponse *cosmostypes.TxResponse
	if req.Mode == tx.BroadcastMode_BROADCAST_MODE_BLOCK {
		txCommitResult, err := c.jsonRPC.BroadcastTxCommit(ctx, req.TxBytes)
		if err != nil {
			return nil, err
		}
		txResponse = cosmostypes.NewResponseFormatBroadcastTxCommit(txCommitResult)
	} else {
		result, err := c.jsonRPC.BroadcastTX(ctx, method, req.TxBytes)
		if err != nil {
			return nil, err
		}
		txResponse = cosmostypes.NewResponseFormatBroadcastTx(result)
	}
	return &tx.BroadcastTxResponse{TxResponse: txResponse}, nil
}

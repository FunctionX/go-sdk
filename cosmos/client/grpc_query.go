package client

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	gogogrpc "github.com/cosmos/gogoproto/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"

	"github.com/functionx/go-sdk/cosmos/codec"
	"github.com/functionx/go-sdk/cosmos/types"
	"github.com/functionx/go-sdk/cosmos/types/tx"
	"github.com/functionx/go-sdk/tendermint/abci"
)

const (
	// BlockHeightHeader is the gRPC header for block height.
	BlockHeightHeader = "x-cosmos-block-height"
)

var _ gogogrpc.ClientConn = &Client{}

// Invoke implements the grpc ClientConn.Invoke method
func (c *Client) Invoke(ctx context.Context, method string, req, reply interface{}, opts ...grpc.CallOption) (err error) {
	// Two things can happen here:
	// 1. either we're broadcasting a Tx, in which call we call CometBFT's broadcast endpoint directly,
	// 2-1. or we are querying for state, in which case we call grpc if grpc client set.
	// 2-2. or we are querying for state, in which case we call ABCI's Query if grpc client not set.

	// In both cases, we don't allow empty request args (it will panic unexpectedly).
	if reflect.ValueOf(req).IsNil() {
		return errors.New("request cannot be nil")
	}

	// Case 1. Broadcasting a Tx.
	if reqProto, ok := req.(*tx.BroadcastTxRequest); ok {
		res, ok := reply.(*tx.BroadcastTxResponse)
		if !ok {
			return fmt.Errorf("expected %T, got %T", (*tx.BroadcastTxResponse)(nil), req)
		}

		broadcastRes, err := c.TxServiceBroadcast(ctx, reqProto)
		if err != nil {
			return err
		}
		*res = *broadcastRes
		return err
	}

	// Case 2-2. Querying state via abci query.
	reqBz, err := c.gRPCCodec().Marshal(req)
	if err != nil {
		return fmt.Errorf("Client.Invoke: failed to marshal request: %w", err)
	}

	// parse height header
	md, _ := metadata.FromOutgoingContext(ctx)
	if heights := md.Get(BlockHeightHeader); len(heights) > 0 {
		height, err := strconv.ParseInt(heights[0], 10, 64)
		if err != nil {
			return err
		}
		if height < 0 {
			return fmt.Errorf("Client.Invoke: height (%d) from %s must be >= 0", height, BlockHeightHeader)
		}

		c.height = height
	}

	res, err := c.QueryABCI(ctx, abci.RequestQuery{Path: method, Data: reqBz, Height: c.height})
	if err != nil {
		return err
	}

	if err = c.gRPCCodec().Unmarshal(res.Value, reply); err != nil {
		return err
	}

	// Create header metadata. For now the headers contain:
	// - block height
	// We then parse all the call options, if the call option is a
	// HeaderCallOption, then we manually set the value of that header to the
	// metadata.
	md = metadata.Pairs(BlockHeightHeader, strconv.FormatInt(res.Height, 10))
	for _, callOpt := range opts {
		header, ok := callOpt.(grpc.HeaderCallOption)
		if !ok {
			continue
		}
		*header.HeaderAddr = md
	}

	if c.cdc.InterfaceRegistry() != nil {
		return types.UnpackInterfaces(reply, c.cdc.InterfaceRegistry())
	}
	return nil
}

// NewStream implements the grpc ClientConn.NewStream method
func (c *Client) NewStream(context.Context, *grpc.StreamDesc, string, ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, fmt.Errorf("streaming rpc not supported")
}

func (c *Client) Close() error {
	return c.jsonRPC.Close()
}

func (c *Client) gRPCCodec() encoding.Codec {
	pc := c.cdc.(codec.GRPCCodecProvider)
	return pc.GRPCCodec()
}

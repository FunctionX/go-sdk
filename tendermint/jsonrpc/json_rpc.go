package jsonrpc

import (
	"context"
	"errors"
	"strconv"

	"github.com/functionx/go-sdk/tendermint/abci"
	"github.com/functionx/go-sdk/tendermint/bytes"
)

type Caller interface {
	Call(ctx context.Context, method string, params map[string]interface{}, result interface{}) error
	Close() error
}

type JsonRPC struct {
	caller Caller
}

func NewJsonRPC(url string) *JsonRPC {
	return &JsonRPC{caller: NewClient(url)}
}

func (c *JsonRPC) ABCIQuery(ctx context.Context, path string, data bytes.HexBytes, height int64, prove bool) (*abci.ResultABCIQuery, error) {
	result := new(abci.ResultABCIQuery)
	params := map[string]interface{}{"path": path, "data": data.String(), "height": strconv.FormatInt(height, 10), "prove": prove}
	err := c.caller.Call(ctx, "abci_query", params, result)
	if err != nil {
		return nil, errors.New("abci query failed: " + err.Error())
	}
	return result, nil
}

func (c *JsonRPC) BroadcastTX(ctx context.Context, method string, tx []byte) (*abci.ResultBroadcastTx, error) {
	result := new(abci.ResultBroadcastTx)
	err := c.caller.Call(ctx, method, map[string]interface{}{"tx": tx}, result)
	if err != nil {
		return nil, errors.New("broadcast tx failed: " + err.Error())
	}
	return result, nil
}

func (c *JsonRPC) BroadcastTxCommit(ctx context.Context, tx []byte) (*abci.ResultBroadcastTxCommit, error) {
	result := new(abci.ResultBroadcastTxCommit)
	err := c.caller.Call(ctx, "broadcast_tx_commit", map[string]interface{}{"tx": tx}, result)
	if err != nil {
		return nil, errors.New("broadcast tx failed: " + err.Error())
	}
	return result, nil
}

func (c *JsonRPC) Close() error {
	return c.caller.Close()
}

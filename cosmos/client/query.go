package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/functionx/go-sdk/tendermint/abci"
)

// QueryABCI performs a query to a CometBFT node with the provide RequestQuery.
// It returns the ResultQuery obtained from the query. The height used to perform
// the query is the RequestQuery Height if it is non-zero, otherwise the context
// height is used.
func (c *Client) QueryABCI(ctx context.Context, req abci.RequestQuery) (abci.ResponseQuery, error) {
	var queryHeight int64
	if req.Height != 0 {
		queryHeight = req.Height
	} else {
		// fallback on the context height
		queryHeight = c.height
	}

	result, err := c.jsonRPC.ABCIQuery(ctx, req.Path, req.Data, queryHeight, req.Prove)
	if err != nil {
		return abci.ResponseQuery{}, err
	}

	if result.Response.Code != 0 {
		return abci.ResponseQuery{}, fmt.Errorf("code: %d, log: %s", result.Response.Code, result.Response.Log)
	}

	// data from trusted node or subspace query doesn't need verification
	if !req.Prove || !isQueryStoreWithProof(req.Path) {
		return result.Response, nil
	}

	return result.Response, nil
}

// QueryStore performs a query to a CometBFT node with the provided key and
// store name. It returns the result and height of the query upon success
// or an error if the query fails.
func (c *Client) QueryStore(key []byte, storeName string) ([]byte, int64, error) {
	return c.queryStore(key, storeName, "key")
}

// queryStore performs a query to a CometBFT node with the provided a store
// name and path. It returns the result and height of the query upon success
// or an error if the query fails.
func (c *Client) queryStore(key []byte, storeName, endPath string) ([]byte, int64, error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, endPath)
	return c.query(path, key)
}

// query performs a query to a CometBFT node with the provided store name
// and path. It returns the result and height of the query upon success
// or an error if the query fails.
func (c *Client) query(path string, key []byte) ([]byte, int64, error) {
	resp, err := c.QueryABCI(context.Background(), abci.RequestQuery{Path: path, Data: key, Height: c.height})
	if err != nil {
		return nil, 0, err
	}
	return resp.Value, resp.Height, nil
}

// isQueryStoreWithProof expects a format like /<queryType>/<storeName>/<subpath>
// queryType must be "store" and subpath must be "key" to require a proof.
func isQueryStoreWithProof(path string) bool {
	if !strings.HasPrefix(path, "/") {
		return false
	}
	paths := strings.SplitN(path[1:], "/", 3)
	switch {
	case len(paths) != 3:
		return false
	case paths[0] != "store":
		return false
	case paths[2] == "key":
		return true
	}
	return false
}

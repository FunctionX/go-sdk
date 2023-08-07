package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/functionx/go-sdk/rand"
)

type RPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      string          `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"` // must be map[string]interface{} or []interface{}
}

type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      string          `json:"id,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func (err RPCError) String() string {
	return fmt.Sprintf("code: %d, data: %s, msg: %s", err.Code, err.Data, err.Message)
}

func NewRPCRequest(id, method string, params json.RawMessage) RPCRequest {
	return RPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}
}

var _ Caller = &Client{}

// Client implement caller
type Client struct {
	Remote string
	client *http.Client
}

func NewClient(remote string) *Client {
	client := http.DefaultClient
	client.Timeout = 10 * time.Second
	return &Client{
		Remote: remote,
		client: client,
	}
}

func (c *Client) Call(ctx context.Context, method string, params map[string]interface{}, result interface{}) (err error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return
	}

	reqId := fmt.Sprintf("go-%s", rand.Str(8))
	body, err := json.Marshal(NewRPCRequest(reqId, method, payload))
	if err != nil {
		return
	}

	if method == "subscribe" {
		return errors.New("this method is not supported")
	}

	if strings.HasPrefix(c.Remote, "tcp") {
		c.Remote = strings.Replace(c.Remote, "tcp", "http", 1)
	}

	// fmt.Println("Request Post ==>", "remote", client.Remote, "body", string(body))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Remote, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	date, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// fmt.Println("Response Body <==", "remote", client.Remote, "body", string(date))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(date))
	}

	var rpcResp RPCResponse
	if err = json.Unmarshal(date, &rpcResp); err != nil {
		return err
	}
	if rpcResp.Error != nil {
		return fmt.Errorf("response code: %d, data: %s, msg: %s", rpcResp.Error.Code, rpcResp.Error.Data, rpcResp.Error.Message)
	}
	return json.Unmarshal(rpcResp.Result, result)
}

func (c *Client) Close() error {
	c.client.CloseIdleConnections()
	return nil
}

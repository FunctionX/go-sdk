package abci

import (
	bytes2 "bytes"

	"github.com/cosmos/gogoproto/jsonpb"

	"github.com/functionx/go-sdk/tendermint/bytes"
)

type ResultABCIQuery struct {
	Response ResponseQuery `json:"response"`
}

type ResultBroadcastTx struct {
	Code      uint32         `json:"code"`
	Data      bytes.HexBytes `json:"data"`
	Log       string         `json:"log"`
	Codespace string         `json:"codespace"`

	Hash bytes.HexBytes `json:"hash"`
}

type ResultBroadcastTxCommit struct {
	CheckTx   ResponseCheckTx   `json:"check_tx"`
	DeliverTx ResponseDeliverTx `json:"deliver_tx"`
	Hash      bytes.HexBytes    `json:"hash"`
	Height    int64             `json:"height"`
}

func (r *ResponseSetOption) UnmarshalJSON(b []byte) error {
	return jsonpb.Unmarshal(bytes2.NewBuffer(b), r)
}

func (r *ResponseCheckTx) UnmarshalJSON(b []byte) error {
	return jsonpb.Unmarshal(bytes2.NewBuffer(b), r)
}

func (r *ResponseDeliverTx) UnmarshalJSON(b []byte) error {
	return jsonpb.Unmarshal(bytes2.NewBuffer(b), r)
}

func (r *ResponseQuery) UnmarshalJSON(b []byte) error {
	return jsonpb.Unmarshal(bytes2.NewBuffer(b), r)
}

func (r *EventAttribute) UnmarshalJSON(b []byte) error {
	return jsonpb.Unmarshal(bytes2.NewBuffer(b), r)
}

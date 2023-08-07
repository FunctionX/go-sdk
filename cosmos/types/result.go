package types

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/functionx/go-sdk/tendermint/abci"
)

func (g *GasInfo) String() string {
	bz, _ := json.Marshal(g)
	return string(bz)
}

func (r *Result) String() string {
	bz, _ := json.Marshal(r)
	return string(bz)
}

func (r *TxResponse) String() string {
	bz, _ := json.Marshal(r)
	return string(bz)
}

type ABCIMessageLogs []ABCIMessageLog

func (logs ABCIMessageLogs) String() string {
	if logs != nil {
		raw, _ := json.Marshal(logs)
		return string(raw)
	}
	return ""
}

func NewResponseFormatBroadcastTx(res *abci.ResultBroadcastTx) *TxResponse {
	if res == nil {
		return nil
	}

	parsedLogs, _ := ParseABCILogs(res.Log)

	return &TxResponse{
		Code:      res.Code,
		Codespace: res.Codespace,
		Data:      res.Data.String(),
		RawLog:    res.Log,
		Logs:      parsedLogs,
		TxHash:    res.Hash.String(),
	}
}

func NewResponseFormatBroadcastTxCommit(res *abci.ResultBroadcastTxCommit) *TxResponse {
	if res == nil {
		return nil
	}

	if res.CheckTx.Code != 0 {
		return newTxResponseCheckTx(res)
	}
	return newTxResponseDeliverTx(res)
}

func newTxResponseCheckTx(res *abci.ResultBroadcastTxCommit) *TxResponse {
	if res == nil {
		return nil
	}
	var txHash string
	if res.Hash != nil {
		txHash = res.Hash.String()
	}

	parsedLogs, _ := ParseABCILogs(res.CheckTx.Log)

	return &TxResponse{
		Height:    res.Height,
		TxHash:    txHash,
		Codespace: res.CheckTx.Codespace,
		Code:      res.CheckTx.Code,
		Data:      strings.ToUpper(hex.EncodeToString(res.CheckTx.Data)),
		RawLog:    res.CheckTx.Log,
		Logs:      parsedLogs,
		Info:      res.CheckTx.Info,
		GasWanted: res.CheckTx.GasWanted,
		GasUsed:   res.CheckTx.GasUsed,
		Events:    res.CheckTx.Events,
	}
}

func newTxResponseDeliverTx(res *abci.ResultBroadcastTxCommit) *TxResponse {
	if res == nil {
		return nil
	}

	var txHash string
	if res.Hash != nil {
		txHash = res.Hash.String()
	}

	parsedLogs, _ := ParseABCILogs(res.DeliverTx.Log)

	return &TxResponse{
		Height:    res.Height,
		TxHash:    txHash,
		Codespace: res.DeliverTx.Codespace,
		Code:      res.DeliverTx.Code,
		Data:      strings.ToUpper(hex.EncodeToString(res.DeliverTx.Data)),
		RawLog:    res.DeliverTx.Log,
		Logs:      parsedLogs,
		Info:      res.DeliverTx.Info,
		GasWanted: res.DeliverTx.GasWanted,
		GasUsed:   res.DeliverTx.GasUsed,
		Events:    res.DeliverTx.Events,
	}
}

func ParseABCILogs(logs string) (res ABCIMessageLogs, err error) {
	return res, json.Unmarshal([]byte(logs), &res)
}

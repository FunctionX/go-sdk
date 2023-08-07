package watcher

import (
	"time"

	"github.com/functionx/go-sdk/cosmos/types/tx"
	"github.com/functionx/go-sdk/tendermint/bytes"
)

type Block struct {
	ChainID   string           `json:"chain_id"`
	Height    int64            `json:"height"`
	BlockTime time.Time        `json:"block_time"`
	TxsHash   []bytes.HexBytes `json:"txs_hash"`
	Txs       []*tx.Tx         `json:"txs"`
}

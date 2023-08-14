package types

import (
	"github.com/cosmos/gogoproto/proto"
)

type (
	// Msg defines the interface a transaction message must fulfill.
	Msg interface {
		proto.Message
		GetSigners() []AccAddress
	}

	Tx interface {
		// GetMsgs returns the all the transaction's messages.
		GetMsgs() []Msg
	}
)

// MsgTypeURL returns the TypeURL of a `sdk.Msg`.
func MsgTypeURL(msg Msg) string {
	return "/" + proto.MessageName(msg)
}

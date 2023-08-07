package transfer

import (
	"fmt"

	"github.com/functionx/go-sdk/cosmos/types"
)

var _ types.Msg = &MsgTransfer{}

func NewMsgTransfer(
	sourcePort,
	sourceChannel string,
	token types.Coin,
	sender,
	receiver string,
	timeoutHeight Height,
	timeoutTimestamp uint64,
	memo string,
) *MsgTransfer {
	return &MsgTransfer{
		SourcePort:       sourcePort,
		SourceChannel:    sourceChannel,
		Token:            token,
		Sender:           sender,
		Receiver:         receiver,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: timeoutTimestamp,
		Memo:             memo,
	}
}

func (msg *MsgTransfer) GetSigners() []types.AccAddress {
	return []types.AccAddress{types.MustAccAddressFromBech32(msg.Sender)}
}

func (h *Height) String() string {
	return fmt.Sprintf("%d-%d", h.RevisionNumber, h.RevisionHeight)
}

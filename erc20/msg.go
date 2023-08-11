package erc20

import (
	"github.com/functionx/go-sdk/cosmos/types"
)

var _ types.Msg = &MsgConvertERC20{}

func NewMsgConvertERC20(
	contractAddress, receiver, sender string, amount types.Int,
) *MsgConvertERC20 {
	return &MsgConvertERC20{
		ContractAddress: contractAddress,
		Amount:          amount,
		Receiver:        receiver,
		Sender:          sender,
	}
}

func (msg *MsgConvertERC20) GetSigners() []types.AccAddress {
	return []types.AccAddress{types.MustAccAddressFromBech32(msg.Sender)}
}

func NewMsgConvertDenom(
	sender, receiver, target string, amount types.Coin,
) *MsgConvertDenom {
	return &MsgConvertDenom{
		Sender:   sender,
		Receiver: receiver,
		Coin:     amount,
		Target:   target,
	}
}

func (msg *MsgConvertDenom) GetSigners() []types.AccAddress {
	return []types.AccAddress{types.MustAccAddressFromBech32(msg.Sender)}
}

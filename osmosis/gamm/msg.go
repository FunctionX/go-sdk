package gamm

import "github.com/functionx/go-sdk/cosmos/types"

var _ types.Msg = (*MsgExitPool)(nil)

func NewMsgExitPool(sender string, poolId uint64, shareInAmount types.Int, outMins types.Coins) *MsgExitPool {
	return &MsgExitPool{
		Sender:        sender,
		PoolId:        poolId,
		ShareInAmount: shareInAmount,
		TokenOutMins:  outMins,
	}
}

func (msg MsgExitPool) GetSigners() []types.AccAddress {
	sender, err := types.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []types.AccAddress{sender}
}

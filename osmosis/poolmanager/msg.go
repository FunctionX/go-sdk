package poolmanager

import "github.com/functionx/go-sdk/cosmos/types"

var _ types.Msg = (*MsgSwapExactAmountIn)(nil)

func NewMsgSwapExactAmountIn(sender string, routers []SwapAmountInRoute, tokenIn types.Coin, minAmountOut types.Int) *MsgSwapExactAmountIn {
	return &MsgSwapExactAmountIn{
		Sender:            sender,
		Routes:            routers,
		TokenIn:           tokenIn,
		TokenOutMinAmount: minAmountOut,
	}
}

func (msg MsgSwapExactAmountIn) GetSigners() []types.AccAddress {
	sender, err := types.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []types.AccAddress{sender}
}

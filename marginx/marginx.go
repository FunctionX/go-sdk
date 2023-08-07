package marginx

import (
	"github.com/functionx/go-sdk/cosmos/types"
)

var (
	_ types.Msg = &MsgCreateOrder{}
	_ types.Msg = &MsgCancelOrder{}
)

func (msg *MsgCreateOrder) GetSigners() []types.AccAddress {
	return []types.AccAddress{types.MustAccAddressFromBech32(msg.Owner)}
}

func (msg *MsgCancelOrder) GetSigners() []types.AccAddress {
	return []types.AccAddress{types.MustAccAddressFromBech32(msg.Owner)}
}

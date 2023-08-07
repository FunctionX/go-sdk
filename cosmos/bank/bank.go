package bank

import (
	"github.com/functionx/go-sdk/cosmos/types"
)

var (
	_ types.Msg = &MsgSend{}
	_ types.Msg = &MsgMultiSend{}
)

func (msg *MsgSend) GetSigners() []types.AccAddress {
	return []types.AccAddress{types.MustAccAddressFromBech32(msg.FromAddress)}
}

func (msg *MsgMultiSend) GetSigners() []types.AccAddress {
	signers := make([]types.AccAddress, len(msg.Inputs))
	for i, input := range msg.Inputs {
		signers[i] = types.MustAccAddressFromBech32(input.Address)
	}
	return signers
}

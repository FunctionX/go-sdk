package erc20

import (
	"github.com/functionx/go-sdk/cosmos/types"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterInterface("cosmos.base.v1beta1.Msg", (*types.Msg)(nil))
	registry.RegisterImplementations((*types.Msg)(nil), &MsgConvertERC20{})
	registry.RegisterImplementations((*types.Msg)(nil), &MsgConvertDenom{})
}

package marginx

import (
	"github.com/functionx/go-sdk/cosmos/codec"
	"github.com/functionx/go-sdk/cosmos/types"
)

func NewProtoCodec() *codec.ProtoCodec {
	registry := types.NewInterfaceRegistry()
	RegisterInterfaces(registry)
	return codec.NewProtoCodec(registry)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterInterface("cosmos.base.v1beta1.Msg", (*types.Msg)(nil))
	registry.RegisterImplementations((*types.Msg)(nil), &MsgCreateOrder{}, &MsgCancelOrder{})
}

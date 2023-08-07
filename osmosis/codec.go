package osmosis

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
}

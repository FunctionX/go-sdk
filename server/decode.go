package server

import (
	"math/big"
	"reflect"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"

	cosmoscrypto "github.com/functionx/go-sdk/cosmos/crypto"
	"github.com/functionx/go-sdk/cosmos/types"
	"github.com/functionx/go-sdk/cosmos/types/tx"
)

func NewConfigDecoder(output interface{}) *mapstructure.Decoder {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: output,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			CustomDecodeHook(),
		),
	})
	if err != nil {
		panic(err)
	}
	return decoder
}

// CustomDecodeHook is a custom decode hook for mapstructure
//
//gocyclo:ignore
func CustomDecodeHook() mapstructure.DecodeHookFunc {
	return mapstructure.ComposeDecodeHookFunc(
		func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			if f.Kind() != reflect.String {
				return data, nil
			}
			if t != reflect.TypeOf(types.Coin{}) {
				return data, nil
			}
			return types.ParseCoin(data.(string))
		},
		func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			if f.Kind() != reflect.String {
				return data, nil
			}
			if t != reflect.TypeOf(big.NewInt(0)) && t != reflect.TypeOf(big.Int{}) {
				return data, nil
			}
			i, err := strconv.ParseInt(data.(string), 10, 64)
			if err != nil {
				return nil, err
			}
			return big.NewInt(i), nil
		},
		func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			if f.Kind() != reflect.String {
				return data, nil
			}
			if t != reflect.TypeOf(tx.BroadcastMode(0)) {
				return data, nil
			}
			return tx.BroadcastMode_value[data.(string)], nil
		},
		func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			if f.Kind() != reflect.String {
				return data, nil
			}
			if t != reflect.TypeOf(cosmoscrypto.Algorithm("")) {
				return data, nil
			}
			return cosmoscrypto.Algorithm(data.(string)), nil
		},
		func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			if f.Kind() != reflect.String {
				return data, nil
			}
			if t != reflect.TypeOf(common.Address{}) {
				return data, nil
			}
			return common.HexToAddress(data.(string)), nil
		},
		func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			if f.Kind() != reflect.String {
				return data, nil
			}
			if t != reflect.TypeOf(types.DecCoin{}) {
				return data, nil
			}
			return types.ParseDecCoin(data.(string))
		},
		func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			if f.Kind() != reflect.Float64 {
				return data, nil
			}
			if t != reflect.TypeOf(decimal.Decimal{}) {
				return data, nil
			}
			return decimal.NewFromFloat(data.(float64)), nil
		},
	)
}

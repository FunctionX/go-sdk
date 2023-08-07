package types

import (
	sdkmath "cosmossdk.io/math"
)

type (
	Int  = sdkmath.Int
	Uint = sdkmath.Uint
	Dec  = sdkmath.LegacyDec
)

func (ip *IntProto) String() string {
	return ip.Int.String()
}

func (dp *DecProto) String() string {
	return dp.Dec.String()
}

func NewDecFromInt(i Int) Dec {
	return sdkmath.LegacyNewDecFromInt(i)
}

func MinDec(a, b Dec) Dec {
	return sdkmath.LegacyMinDec(a, b)
}

func ZeroDec() Dec {
	return sdkmath.LegacyZeroDec()
}

func NewDecFromStr(s string) (Dec, error) {
	return sdkmath.LegacyNewDecFromStr(s)
}

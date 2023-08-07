package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	sdkmath "cosmossdk.io/math"
)

func NewCoin(denom string, amount Int) Coin {
	coin := Coin{
		Denom:  denom,
		Amount: amount,
	}

	if err := coin.Validate(); err != nil {
		panic(err)
	}

	return coin
}

func NewInt64Coin(denom string, amount int64) Coin {
	return NewCoin(denom, sdkmath.NewInt(amount))
}

func (coin Coin) String() string {
	return fmt.Sprintf("%v%s", coin.Amount, coin.Denom)
}

func (coin Coin) Validate() error {
	if err := ValidateDenom(coin.Denom); err != nil {
		return err
	}

	if coin.Amount.IsNegative() {
		return fmt.Errorf("negative coin amount: %v", coin.Amount)
	}

	return nil
}

func (coin Coin) IsValid() bool {
	return coin.Validate() == nil
}

func (coin Coin) IsZero() bool {
	return coin.Amount.IsZero()
}

func (coin Coin) IsGTE(other Coin) bool {
	if coin.Denom != other.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", coin.Denom, other.Denom))
	}

	return !coin.Amount.LT(other.Amount)
}

func (coin Coin) IsLT(other Coin) bool {
	if coin.Denom != other.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", coin.Denom, other.Denom))
	}

	return coin.Amount.LT(other.Amount)
}

func (coin Coin) IsLTE(other Coin) bool {
	if coin.Denom != other.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", coin.Denom, other.Denom))
	}

	return !coin.Amount.GT(other.Amount)
}

func (coin Coin) IsEqual(other Coin) bool {
	if coin.Denom != other.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", coin.Denom, other.Denom))
	}

	return coin.Amount.Equal(other.Amount)
}

func (coin Coin) Add(coinB Coin) Coin {
	if coin.Denom != coinB.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", coin.Denom, coinB.Denom))
	}

	return Coin{coin.Denom, coin.Amount.Add(coinB.Amount)}
}

func (coin Coin) AddAmount(amount Int) Coin {
	return Coin{coin.Denom, coin.Amount.Add(amount)}
}

func (coin Coin) Sub(coinB Coin) Coin {
	res, err := coin.SafeSub(coinB)
	if err != nil {
		panic(err)
	}

	return res
}

func (coin Coin) SafeSub(coinB Coin) (Coin, error) {
	if coin.Denom != coinB.Denom {
		return Coin{}, fmt.Errorf("invalid coin denoms: %s, %s", coin.Denom, coinB.Denom)
	}

	res := Coin{coin.Denom, coin.Amount.Sub(coinB.Amount)}
	if res.IsNegative() {
		return Coin{}, fmt.Errorf("negative coin amount")
	}

	return res, nil
}

func (coin Coin) SubAmount(amount Int) Coin {
	res := Coin{coin.Denom, coin.Amount.Sub(amount)}
	if res.IsNegative() {
		panic("negative coin amount")
	}

	return res
}

func (coin Coin) IsPositive() bool {
	return coin.Amount.Sign() == 1
}

func (coin Coin) IsNegative() bool {
	return coin.Amount.Sign() == -1
}

func (coin Coin) IsNil() bool {
	return coin.Amount.BigInt() == nil
}

type Coins []Coin

func NewCoins(coins ...Coin) Coins {
	newCoins := sanitizeCoins(coins)
	if err := newCoins.Validate(); err != nil {
		panic(fmt.Errorf("invalid coin set %s: %w", newCoins, err))
	}

	return newCoins
}

func sanitizeCoins(coins []Coin) Coins {
	newCoins := removeZeroCoins(coins)
	if len(newCoins) == 0 {
		return Coins{}
	}

	return newCoins.Sort()
}

type coinsJSON Coins

func (coins Coins) MarshalJSON() ([]byte, error) {
	if coins == nil {
		return json.Marshal(coinsJSON(Coins{}))
	}

	return json.Marshal(coinsJSON(coins))
}

func (coins Coins) String() string {
	if len(coins) == 0 {
		return ""
	} else if len(coins) == 1 {
		return coins[0].String()
	}

	var out strings.Builder
	for _, coin := range coins[:len(coins)-1] {
		out.WriteString(coin.String())
		out.WriteByte(',')
	}
	out.WriteString(coins[len(coins)-1].String())
	return out.String()
}

func (coins Coins) Validate() error {
	switch len(coins) {
	case 0:
		return nil

	case 1:
		if err := ValidateDenom(coins[0].Denom); err != nil {
			return err
		}
		if !coins[0].IsPositive() {
			return fmt.Errorf("coin %s amount is not positive", coins[0])
		}
		return nil

	default:

		if err := (Coins{coins[0]}).Validate(); err != nil {
			return err
		}

		lowDenom := coins[0].Denom
		seenDenoms := make(map[string]bool)
		seenDenoms[lowDenom] = true

		for _, coin := range coins[1:] {
			if seenDenoms[coin.Denom] {
				return fmt.Errorf("duplicate denomination %s", coin.Denom)
			}
			if err := ValidateDenom(coin.Denom); err != nil {
				return err
			}
			if coin.Denom <= lowDenom {
				return fmt.Errorf("denomination %s is not sorted", coin.Denom)
			}
			if !coin.IsPositive() {
				return fmt.Errorf("coin %s amount is not positive", coin.Denom)
			}

			lowDenom = coin.Denom
			seenDenoms[coin.Denom] = true
		}

		return nil
	}
}

func (coins Coins) isSorted() bool {
	for i := 1; i < len(coins); i++ {
		if coins[i-1].Denom > coins[i].Denom {
			return false
		}
	}
	return true
}

func (coins Coins) IsValid() bool {
	return coins.Validate() == nil
}

func (coins Coins) Add(coinsB ...Coin) Coins {
	return coins.safeAdd(coinsB)
}

func (coins Coins) safeAdd(coinsB Coins) (coalesced Coins) {
	if !coins.isSorted() {
		panic("Coins (self) must be sorted")
	}
	if !coinsB.isSorted() {
		panic("Wrong argument: coins must be sorted")
	}

	uniqCoins := make(map[string]Coins, len(coins)+len(coinsB))

	for _, cL := range []Coins{coins, coinsB} {
		for _, c := range cL {
			uniqCoins[c.Denom] = append(uniqCoins[c.Denom], c)
		}
	}

	for denom, cL := range uniqCoins {
		comboCoin := Coin{Denom: denom, Amount: sdkmath.NewInt(0)}
		for _, c := range cL {
			comboCoin = comboCoin.Add(c)
		}
		if !comboCoin.IsZero() {
			coalesced = append(coalesced, comboCoin)
		}
	}
	return coalesced.Sort()
}

func (coins Coins) DenomsSubsetOf(coinsB Coins) bool {
	if len(coins) > len(coinsB) {
		return false
	}

	for _, coin := range coins {
		if coinsB.AmountOf(coin.Denom).IsZero() {
			return false
		}
	}

	return true
}

func (coins Coins) Sub(coinsB ...Coin) Coins {
	diff, hasNeg := coins.SafeSub(coinsB...)
	if hasNeg {
		panic("negative coin amount")
	}

	return diff
}

func (coins Coins) SafeSub(coinsB ...Coin) (Coins, bool) {
	diff := coins.safeAdd(NewCoins(coinsB...).negative())
	return diff, diff.IsAnyNegative()
}

func (coins Coins) MulInt(x Int) Coins {
	coins, ok := coins.SafeMulInt(x)
	if !ok {
		panic("multiplying by zero is an invalid operation on coins")
	}

	return coins
}

func (coins Coins) SafeMulInt(x Int) (Coins, bool) {
	if x.IsZero() {
		return nil, false
	}

	res := make(Coins, len(coins))
	for i, coin := range coins {
		coin := coin
		res[i] = NewCoin(coin.Denom, coin.Amount.Mul(x))
	}

	return res, true
}

func (coins Coins) QuoInt(x Int) Coins {
	coins, ok := coins.SafeQuoInt(x)
	if !ok {
		panic("dividing by zero is an invalid operation on coins")
	}

	return coins
}

func (coins Coins) SafeQuoInt(x Int) (Coins, bool) {
	if x.IsZero() {
		return nil, false
	}

	res := make(Coins, len(coins))
	for _, coin := range coins {
		res = append(res, NewCoin(coin.Denom, coin.Amount.Quo(x)))
	}

	return res, true
}

func (coins Coins) Max(coinsB Coins) Coins {
	max := make([]Coin, 0)
	indexA, indexB := 0, 0
	for indexA < len(coins) && indexB < len(coinsB) {
		coinA, coinB := coins[indexA], coinsB[indexB]
		switch strings.Compare(coinA.Denom, coinB.Denom) {
		case -1:
			max = append(max, coinA)
			indexA++
		case 0:
			maxCoin := coinA
			if coinB.Amount.GT(maxCoin.Amount) {
				maxCoin = coinB
			}
			max = append(max, maxCoin)
			indexA++
			indexB++
		case 1:
			max = append(max, coinB)
			indexB++
		}
	}
	for ; indexA < len(coins); indexA++ {
		max = append(max, coins[indexA])
	}
	for ; indexB < len(coinsB); indexB++ {
		max = append(max, coinsB[indexB])
	}
	return NewCoins(max...)
}

func (coins Coins) Min(coinsB Coins) Coins {
	min := make([]Coin, 0)
	for indexA, indexB := 0, 0; indexA < len(coins) && indexB < len(coinsB); {
		coinA, coinB := coins[indexA], coinsB[indexB]
		switch strings.Compare(coinA.Denom, coinB.Denom) {
		case -1:
			indexA++
		case 0:
			minCoin := coinA
			if coinB.Amount.LT(minCoin.Amount) {
				minCoin = coinB
			}
			if !minCoin.IsZero() {
				min = append(min, minCoin)
			}
			indexA++
			indexB++
		case 1:
			indexB++
		}
	}
	return NewCoins(min...)
}

func (coins Coins) IsAllGT(coinsB Coins) bool {
	if len(coins) == 0 {
		return false
	}

	if len(coinsB) == 0 {
		return true
	}

	if !coinsB.DenomsSubsetOf(coins) {
		return false
	}

	for _, coinB := range coinsB {
		amountA, amountB := coins.AmountOf(coinB.Denom), coinB.Amount
		if !amountA.GT(amountB) {
			return false
		}
	}

	return true
}

func (coins Coins) IsAllGTE(coinsB Coins) bool {
	if len(coinsB) == 0 {
		return true
	}

	if len(coins) == 0 {
		return false
	}

	for _, coinB := range coinsB {
		if coinB.Amount.GT(coins.AmountOf(coinB.Denom)) {
			return false
		}
	}

	return true
}

func (coins Coins) IsAllLT(coinsB Coins) bool {
	return coinsB.IsAllGT(coins)
}

func (coins Coins) IsAllLTE(coinsB Coins) bool {
	return coinsB.IsAllGTE(coins)
}

func (coins Coins) IsAnyGT(coinsB Coins) bool {
	if len(coinsB) == 0 {
		return false
	}

	for _, coin := range coins {
		amt := coinsB.AmountOf(coin.Denom)
		if coin.Amount.GT(amt) && !amt.IsZero() {
			return true
		}
	}

	return false
}

func (coins Coins) IsAnyGTE(coinsB Coins) bool {
	if len(coinsB) == 0 {
		return false
	}

	for _, coin := range coins {
		amt := coinsB.AmountOf(coin.Denom)
		if coin.Amount.GTE(amt) && !amt.IsZero() {
			return true
		}
	}

	return false
}

func (coins Coins) IsZero() bool {
	for _, coin := range coins {
		if !coin.IsZero() {
			return false
		}
	}
	return true
}

func (coins Coins) IsEqual(coinsB Coins) bool {
	if len(coins) != len(coinsB) {
		return false
	}

	coins = coins.Sort()
	coinsB = coinsB.Sort()

	for i := 0; i < len(coins); i++ {
		if !coins[i].IsEqual(coinsB[i]) {
			return false
		}
	}

	return true
}

func (coins Coins) Empty() bool {
	return len(coins) == 0
}

func (coins Coins) AmountOf(denom string) Int {
	if err := ValidateDenom(denom); err != nil {
		panic(err)
	}
	return coins.AmountOfNoDenomValidation(denom)
}

func (coins Coins) AmountOfNoDenomValidation(denom string) Int {
	if ok, c := coins.Find(denom); ok {
		return c.Amount
	}
	return sdkmath.ZeroInt()
}

func (coins Coins) Find(denom string) (bool, Coin) {
	switch len(coins) {
	case 0:
		return false, Coin{}

	case 1:
		coin := coins[0]
		if coin.Denom == denom {
			return true, coin
		}
		return false, Coin{}

	default:
		midIdx := len(coins) / 2
		coin := coins[midIdx]
		switch {
		case denom < coin.Denom:
			return coins[:midIdx].Find(denom)
		case denom == coin.Denom:
			return true, coin
		default:
			return coins[midIdx+1:].Find(denom)
		}
	}
}

func (coins Coins) GetDenomByIndex(i int) string {
	return coins[i].Denom
}

func (coins Coins) IsAllPositive() bool {
	if len(coins) == 0 {
		return false
	}

	for _, coin := range coins {
		if !coin.IsPositive() {
			return false
		}
	}

	return true
}

func (coins Coins) IsAnyNegative() bool {
	for _, coin := range coins {
		if coin.IsNegative() {
			return true
		}
	}

	return false
}

func (coins Coins) IsAnyNil() bool {
	for _, coin := range coins {
		if coin.IsNil() {
			return true
		}
	}

	return false
}

func (coins Coins) negative() Coins {
	res := make([]Coin, 0, len(coins))

	for _, coin := range coins {
		res = append(res, Coin{
			Denom:  coin.Denom,
			Amount: coin.Amount.Neg(),
		})
	}

	return res
}

func removeZeroCoins(coins Coins) Coins {
	for i := 0; i < len(coins); i++ {
		if coins[i].IsZero() {
			break
		} else if i == len(coins)-1 {
			return coins
		}
	}

	var result []Coin
	if len(coins) > 0 {
		result = make([]Coin, 0, len(coins)-1)
	}

	for _, coin := range coins {
		if !coin.IsZero() {
			result = append(result, coin)
		}
	}

	return result
}

func (coins Coins) Len() int { return len(coins) }

func (coins Coins) Less(i, j int) bool { return coins[i].Denom < coins[j].Denom }

func (coins Coins) Swap(i, j int) { coins[i], coins[j] = coins[j], coins[i] }

var _ sort.Interface = Coins{}

func (coins Coins) Sort() Coins {
	sort.Sort(coins)
	return coins
}

func ParseCoins(coinsStr string) (Coins, error) {
	decCoins, err := ParseDecCoins(coinsStr)
	if err != nil {
		return nil, err
	}
	if decCoins == nil {
		return nil, nil
	}
	result := make([]Coin, 0, len(decCoins))

	for _, coin := range decCoins {
		newCoin, _ := coin.TruncateDecimal()
		result = append(result, newCoin)
	}

	return result, nil
}

func ParseCoin(coinStr string) (Coin, error) {
	decCoin, err := ParseDecCoin(coinStr)
	if err != nil {
		return Coin{}, err
	}
	newCoin, _ := decCoin.TruncateDecimal()
	return newCoin, nil
}

func ValidateDenom(denom string) error {
	if !reDnm.MatchString(denom) {
		return fmt.Errorf("invalid denom: %s", denom)
	}
	return nil
}

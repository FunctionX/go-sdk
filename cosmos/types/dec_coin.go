package types

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/pkg/errors"
)

var (
	reDnmString = `[a-zA-Z][a-zA-Z0-9/:._-]{1,127}`
	reDecAmt    = `[[:digit:]]+(?:\.[[:digit:]]+)?|\.[[:digit:]]+`
	reSpc       = `[[:space:]]*`
	reDnm       = regexp.MustCompile(fmt.Sprintf(`^%s$`, reDnmString))
	reDecCoin   = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reDecAmt, reSpc, reDnmString))
)

func NewDecCoin(denom string, amount Int) DecCoin {
	coin := NewCoin(denom, amount)

	return DecCoin{
		Denom:  coin.Denom,
		Amount: NewDecFromInt(coin.Amount),
	}
}

func NewDecCoinFromDec(denom string, amount Dec) DecCoin {
	if err := ValidateDenom(denom); err != nil {
		panic(err)
	}
	if amount.IsNegative() {
		panic(fmt.Sprintf("negative decimal coin amount: %v\n", amount))
	}

	return DecCoin{
		Denom:  denom,
		Amount: amount,
	}
}

func NewDecCoinFromCoin(coin Coin) DecCoin {
	return DecCoin{
		Denom:  coin.Denom,
		Amount: NewDecFromInt(coin.Amount),
	}
}

func NewInt64DecCoin(denom string, amount int64) DecCoin {
	return NewDecCoin(denom, sdkmath.NewInt(amount))
}

func (coin DecCoin) IsZero() bool {
	return coin.Amount.IsZero()
}

func (coin DecCoin) IsGTE(other DecCoin) bool {
	if coin.Denom != other.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", coin.Denom, other.Denom))
	}

	return !coin.Amount.LT(other.Amount)
}

func (coin DecCoin) IsLT(other DecCoin) bool {
	if coin.Denom != other.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", coin.Denom, other.Denom))
	}

	return coin.Amount.LT(other.Amount)
}

func (coin DecCoin) IsEqual(other DecCoin) bool {
	if coin.Denom != other.Denom {
		panic(fmt.Sprintf("invalid coin denominations; %s, %s", coin.Denom, other.Denom))
	}

	return coin.Amount.Equal(other.Amount)
}

func (coin DecCoin) Add(coinB DecCoin) DecCoin {
	if coin.Denom != coinB.Denom {
		panic(fmt.Sprintf("coin denom different: %v %v\n", coin.Denom, coinB.Denom))
	}
	return DecCoin{coin.Denom, coin.Amount.Add(coinB.Amount)}
}

func (coin DecCoin) Sub(coinB DecCoin) DecCoin {
	if coin.Denom != coinB.Denom {
		panic(fmt.Sprintf("coin denom different: %v %v\n", coin.Denom, coinB.Denom))
	}
	res := DecCoin{coin.Denom, coin.Amount.Sub(coinB.Amount)}
	if res.IsNegative() {
		panic("negative decimal coin amount")
	}
	return res
}

func (coin DecCoin) TruncateDecimal() (Coin, DecCoin) {
	truncated := coin.Amount.TruncateInt()
	change := coin.Amount.Sub(NewDecFromInt(truncated))
	return NewCoin(coin.Denom, truncated), NewDecCoinFromDec(coin.Denom, change)
}

func (coin DecCoin) IsPositive() bool {
	return coin.Amount.IsPositive()
}

func (coin DecCoin) IsNegative() bool {
	return coin.Amount.IsNegative()
}

func (coin DecCoin) String() string {
	return fmt.Sprintf("%v%v", coin.Amount, coin.Denom)
}

func (coin DecCoin) Validate() error {
	if err := ValidateDenom(coin.Denom); err != nil {
		return err
	}
	if coin.IsNegative() {
		return fmt.Errorf("decimal coin %s amount cannot be negative", coin)
	}
	return nil
}

func (coin DecCoin) IsValid() bool {
	return coin.Validate() == nil
}

type DecCoins []DecCoin

func NewDecCoins(decCoins ...DecCoin) DecCoins {
	newDecCoins := sanitizeDecCoins(decCoins)
	if err := newDecCoins.Validate(); err != nil {
		panic(fmt.Errorf("invalid coin set %s: %w", newDecCoins, err))
	}

	return newDecCoins
}

func sanitizeDecCoins(decCoins []DecCoin) DecCoins {
	newDecCoins := removeZeroDecCoins(decCoins)
	if len(newDecCoins) == 0 {
		return DecCoins{}
	}

	return newDecCoins.Sort()
}

func NewDecCoinsFromCoins(coins ...Coin) DecCoins {
	decCoins := make(DecCoins, len(coins))
	newCoins := NewCoins(coins...)
	for i, coin := range newCoins {
		decCoins[i] = NewDecCoinFromCoin(coin)
	}

	return decCoins
}

func (coins DecCoins) String() string {
	if len(coins) == 0 {
		return ""
	}

	out := ""
	for _, coin := range coins {
		out += fmt.Sprintf("%v,", coin.String())
	}

	return out[:len(out)-1]
}

func (coins DecCoins) TruncateDecimal() (truncatedCoins Coins, changeCoins DecCoins) {
	for _, coin := range coins {
		truncated, change := coin.TruncateDecimal()
		if !truncated.IsZero() {
			truncatedCoins = truncatedCoins.Add(truncated)
		}
		if !change.IsZero() {
			changeCoins = changeCoins.Add(change)
		}
	}

	return truncatedCoins, changeCoins
}

func (coins DecCoins) Add(coinsB ...DecCoin) DecCoins {
	return coins.safeAdd(coinsB)
}

func (coins DecCoins) safeAdd(coinsB DecCoins) DecCoins {
	sum := ([]DecCoin)(nil)
	indexA, indexB := 0, 0
	lenA, lenB := len(coins), len(coinsB)

	for {
		if indexA == lenA {
			if indexB == lenB {
				return sum
			}

			return append(sum, removeZeroDecCoins(coinsB[indexB:])...)
		} else if indexB == lenB {
			return append(sum, removeZeroDecCoins(coins[indexA:])...)
		}

		coinA, coinB := coins[indexA], coinsB[indexB]

		switch strings.Compare(coinA.Denom, coinB.Denom) {
		case -1:
			if !coinA.IsZero() {
				sum = append(sum, coinA)
			}

			indexA++

		case 0:
			res := coinA.Add(coinB)
			if !res.IsZero() {
				sum = append(sum, res)
			}

			indexA++
			indexB++

		case 1:
			if !coinB.IsZero() {
				sum = append(sum, coinB)
			}

			indexB++
		}
	}
}

func (coins DecCoins) negative() DecCoins {
	res := make([]DecCoin, 0, len(coins))
	for _, coin := range coins {
		res = append(res, DecCoin{
			Denom:  coin.Denom,
			Amount: coin.Amount.Neg(),
		})
	}
	return res
}

func (coins DecCoins) Sub(coinsB DecCoins) DecCoins {
	diff, hasNeg := coins.SafeSub(coinsB)
	if hasNeg {
		panic("negative coin amount")
	}

	return diff
}

func (coins DecCoins) SafeSub(coinsB DecCoins) (DecCoins, bool) {
	diff := coins.safeAdd(coinsB.negative())
	return diff, diff.IsAnyNegative()
}

func (coins DecCoins) Intersect(coinsB DecCoins) DecCoins {
	res := make([]DecCoin, len(coins))
	for i, coin := range coins {
		minCoin := DecCoin{
			Denom:  coin.Denom,
			Amount: MinDec(coin.Amount, coinsB.AmountOf(coin.Denom)),
		}
		res[i] = minCoin
	}
	return removeZeroDecCoins(res)
}

func (coins DecCoins) GetDenomByIndex(i int) string {
	return coins[i].Denom
}

func (coins DecCoins) IsAnyNegative() bool {
	for _, coin := range coins {
		if coin.IsNegative() {
			return true
		}
	}

	return false
}

func (coins DecCoins) MulDec(d Dec) DecCoins {
	var res DecCoins
	for _, coin := range coins {
		product := DecCoin{
			Denom:  coin.Denom,
			Amount: coin.Amount.Mul(d),
		}

		if !product.IsZero() {
			res = res.Add(product)
		}
	}

	return res
}

func (coins DecCoins) MulDecTruncate(d Dec) DecCoins {
	if d.IsZero() {
		return DecCoins{}
	}

	var res DecCoins
	for _, coin := range coins {
		product := DecCoin{
			Denom:  coin.Denom,
			Amount: coin.Amount.MulTruncate(d),
		}

		if !product.IsZero() {
			res = res.Add(product)
		}
	}

	return res
}

func (coins DecCoins) QuoDec(d Dec) DecCoins {
	if d.IsZero() {
		panic("invalid zero decimal")
	}

	var res DecCoins
	for _, coin := range coins {
		quotient := DecCoin{
			Denom:  coin.Denom,
			Amount: coin.Amount.Quo(d),
		}

		if !quotient.IsZero() {
			res = res.Add(quotient)
		}
	}

	return res
}

func (coins DecCoins) QuoDecTruncate(d Dec) DecCoins {
	if d.IsZero() {
		panic("invalid zero decimal")
	}

	var res DecCoins
	for _, coin := range coins {
		quotient := DecCoin{
			Denom:  coin.Denom,
			Amount: coin.Amount.QuoTruncate(d),
		}

		if !quotient.IsZero() {
			res = res.Add(quotient)
		}
	}

	return res
}

func (coins DecCoins) Empty() bool {
	return len(coins) == 0
}

func (coins DecCoins) AmountOf(denom string) Dec {
	if err := ValidateDenom(denom); err != nil {
		panic(err)
	}
	switch len(coins) {
	case 0:
		return ZeroDec()

	case 1:
		coin := coins[0]
		if coin.Denom == denom {
			return coin.Amount
		}
		return ZeroDec()

	default:
		midIdx := len(coins) / 2
		coin := coins[midIdx]

		switch {
		case denom < coin.Denom:
			return coins[:midIdx].AmountOf(denom)
		case denom == coin.Denom:
			return coin.Amount
		default:
			return coins[midIdx+1:].AmountOf(denom)
		}
	}
}

func (coins DecCoins) IsEqual(coinsB DecCoins) bool {
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

func (coins DecCoins) IsZero() bool {
	for _, coin := range coins {
		if !coin.Amount.IsZero() {
			return false
		}
	}
	return true
}

func (coins DecCoins) Validate() error {
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

		if err := (DecCoins{coins[0]}).Validate(); err != nil {
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

func (coins DecCoins) IsValid() bool {
	return coins.Validate() == nil
}

func (coins DecCoins) IsAllPositive() bool {
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

func removeZeroDecCoins(coins DecCoins) DecCoins {
	result := make([]DecCoin, 0, len(coins))

	for _, coin := range coins {
		if !coin.IsZero() {
			result = append(result, coin)
		}
	}

	return result
}

var _ sort.Interface = DecCoins{}

func (coins DecCoins) Len() int { return len(coins) }

func (coins DecCoins) Less(i, j int) bool { return coins[i].Denom < coins[j].Denom }

func (coins DecCoins) Swap(i, j int) { coins[i], coins[j] = coins[j], coins[i] }

func (coins DecCoins) Sort() DecCoins {
	sort.Sort(coins)
	return coins
}

func ParseDecCoin(coinStr string) (coin DecCoin, err error) {
	coinStr = strings.TrimSpace(coinStr)

	matches := reDecCoin.FindStringSubmatch(coinStr)
	if matches == nil {
		return DecCoin{}, fmt.Errorf("invalid decimal coin expression: %s", coinStr)
	}

	amountStr, denomStr := matches[1], matches[2]

	amount, err := NewDecFromStr(amountStr)
	if err != nil {
		return DecCoin{}, errors.Wrap(err, fmt.Sprintf("failed to parse decimal coin amount: %s", amountStr))
	}

	if err = ValidateDenom(denomStr); err != nil {
		return DecCoin{}, fmt.Errorf("invalid denom cannot contain spaces: %w", err)
	}

	return NewDecCoinFromDec(denomStr, amount), nil
}

func ParseDecCoins(coinsStr string) (DecCoins, error) {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		return nil, nil
	}

	coinStrs := strings.Split(coinsStr, ",")
	decCoins := make(DecCoins, len(coinStrs))
	for i, coinStr := range coinStrs {
		coin, err := ParseDecCoin(coinStr)
		if err != nil {
			return nil, err
		}

		decCoins[i] = coin
	}

	newDecCoins := sanitizeDecCoins(decCoins)
	if err := newDecCoins.Validate(); err != nil {
		return nil, err
	}

	return newDecCoins, nil
}

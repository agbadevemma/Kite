package money

import "math"

var currencyDecimals = map[string]int{
	"USD": 2,
	"EUR": 2,
	"GBP": 2,
	"NGN": 2,
	"KES": 2,
	"JPY": 0,
}

func ToSmallestUnit(amount float64, currency string) int64 {

	decimals := currencyDecimals[currency]

	multiplier := math.Pow10(decimals)

	return int64(math.Round(amount * multiplier))
}


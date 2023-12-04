package util

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func ParseDecimalOrZero(str string) decimal.Decimal {
	d, err := decimal.NewFromString(str)
	if err != nil {
		fmt.Printf("Error parsing string to decimal. Reason: ", err.Error())
		return decimal.NewFromInt(0)
	}
	return d
}

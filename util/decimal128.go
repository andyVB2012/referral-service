package util

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/big"
)

func AddDecimal128(a *primitive.Decimal128, b *primitive.Decimal128) *primitive.Decimal128 {
	aF, _, _ := big.ParseFloat(a.String(), 10, 128, big.ToNearestEven)
	bF, _, _ := big.ParseFloat(b.String(), 10, 128, big.ToNearestEven)
	sum := big.NewFloat(0).Add(aF, bF)

	decimal128, _ := primitive.ParseDecimal128(sum.String())
	return &decimal128
}

func SubDecimal128(a *primitive.Decimal128, b *primitive.Decimal128) *primitive.Decimal128 {
	aF, _, _ := big.ParseFloat(a.String(), 10, 128, big.ToNearestEven)
	bF, _, _ := big.ParseFloat(b.String(), 10, 128, big.ToNearestEven)
	sum := big.NewFloat(0).Sub(aF, bF)

	decimal128, _ := primitive.ParseDecimal128(sum.String())
	return &decimal128
}

func ToBigFloat(a *primitive.Decimal128) *big.Float {
	aF, _, _ := big.ParseFloat(a.String(), 10, 128, big.ToNearestEven)
	return aF
}

func ToFloat64(a *primitive.Decimal128) float64 {
	aF, _, _ := big.ParseFloat(a.String(), 10, 128, big.ToNearestEven)
	f, _ := aF.Float64()
	return f
}

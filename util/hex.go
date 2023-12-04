package util

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

func DecodeBigHexOrPanic(hex string) *big.Int {
	d, err := hexutil.DecodeBig(hex)
	if err != nil {
		panic("error decoding hex string to big int. Reason: " + err.Error())
	}
	return d
}

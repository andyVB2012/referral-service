package util

import (
	"fmt"
	"strconv"
)

func ParseUintOrDefault(str string, base int, bitSize int) uint64 {
	parseUint, err := strconv.ParseUint(str, base, bitSize)
	if err != nil {
		fmt.Printf("Error converting string to int. Reason: ", err.Error())
		return 0
	}
	return parseUint
}

func ParseInt64OrZero(str string) int64 {
	p, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return p
}

func ParseInt64OrDefault(str string, def int64) int64 {
	p, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return def
	}
	return p
}

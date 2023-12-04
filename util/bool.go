package util

import (
	"fmt"
	"strconv"
)

func ParseBoolOrDefault(str string) bool {
	b, err := strconv.ParseBool(str)
	if err != nil {
		fmt.Printf("Error converting string to bool. Reason: ", err.Error())
		return false
	}
	return b
}

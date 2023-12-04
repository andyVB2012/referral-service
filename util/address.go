package util

import "errors"

// ShortenEthAddress - shorten the received eth address
// 0xA3B353A6107a3b0aD4239759AB733F8B5f80034a -> 0xA3B353...0034a
func ShortenEthAddress(address string) (string, error) {
	if IsValidAddress(address) {
		return address[:6] + "..." + address[len(address)-4:], nil
	}
	return address, errors.New("received address is not a valid ethereum address")
}

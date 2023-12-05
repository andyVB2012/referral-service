package http

import (
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func VerifySignature(message, signature, address string) (bool, error) {
	sig := hexutil.MustDecode(signature)
	msg := accounts.TextHash([]byte(message))

	// Ethereum signatures are [R || S || V] format where V is 0 or 1
	if len(sig) != 65 {
		return false, errors.New("invalid signature")
	}

	sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1

	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {

		return false, err
	}

	recoveredAddr := crypto.PubkeyToAddress(*recovered)
	verified := address == strings.ToLower(recoveredAddr.Hex())

	return verified, nil
}

func Error(s string) {
	panic("unimplemented")
}

package util

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"math/big"
)

// MBigInt Mongodb big int
type MBigInt struct {
	BigInt *big.Int
}

func NewMBigInt(bint *big.Int) *MBigInt {
	return &MBigInt{BigInt: bint}
}

func NewMBigIntWithZero() *MBigInt {
	return &MBigInt{big.NewInt(0)}
}

func (i *MBigInt) MarshalBSONValue() (bsontype.Type, []byte, error) {
	var v string
	if i != nil {
		v = i.BigInt.String()
	} else {
		v = "0"
	}
	return bson.MarshalValue(v)
}

func (i *MBigInt) UnmarshalBSONValue(t bsontype.Type, value []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("invalid bson value model '%s'", t.String())
	}
	str, _, ok := bsoncore.ReadString(value)
	if !ok {
		return fmt.Errorf("invalid bson string value")
	}
	bigint, err := ParseBigInt(str)
	if err != nil {
		return fmt.Errorf("error converting string to bigint ('%s')", err)
	}
	i.BigInt = bigint
	return nil
}

func (i *MBigInt) UnmarshalJSON(data []byte) error {
	i.BigInt = ParseBigIntOrDefault(string(data))
	return nil
}

func (i *MBigInt) Gt(b *MBigInt) bool {
	cmp := i.BigInt.Cmp(b.BigInt)
	return cmp == 1
}

func (i *MBigInt) Gte(b *MBigInt) bool {
	cmp := i.BigInt.Cmp(b.BigInt)
	return cmp == 0 || cmp == 1
}

func (i *MBigInt) Eq(b *MBigInt) bool {
	cmp := i.BigInt.Cmp(b.BigInt)
	return cmp == 0
}

func (i *MBigInt) String() string {
	return i.BigInt.String()
}

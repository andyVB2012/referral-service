package block

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/andyVB2012/referral-service/util"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func UnmarshalString(bindex *BlockIndex, key string) string {
	str := wrapperspb.String("")
	err := bindex.EventArgs[key].UnmarshalTo(str)
	if err != nil {
		fmt.Printf("Error unmarshalling string key &s on event %s. Reason: %s", key, bindex.Event, err.Error())
		return ""
	}
	return str.Value
}

func UnmarshalBool(bindex *BlockIndex, key string) bool {
	b := wrapperspb.Bool(false)
	err := bindex.EventArgs[key].UnmarshalTo(b)
	if err != nil {
		fmt.Printf("Error unmarshalling bool key &s on event %s. Reason: %s", key, bindex.Event, err.Error())
		return false
	}
	return b.Value
}

func UnmarshallArrMBigInt(bindex *BlockIndex, key string) []*util.MBigInt {
	str := wrapperspb.String("")
	if err := bindex.EventArgs[key].UnmarshalTo(str); err != nil {
		fmt.Printf("Error unmarshalling arr big int key &s on event %s. Reason: %s", key, bindex.Event, err.Error())
		return []*util.MBigInt{}
	}

	strArrConv := strings.ReplaceAll(str.Value, " ", ",")

	var arr []*util.MBigInt
	if err := json.Unmarshal([]byte(strArrConv), &arr); err != nil {
		fmt.Printf("Error unmarshalling arr big int key &s on event %s. Reason: %s", key, bindex.Event, err.Error())
		return []*util.MBigInt{}
	}

	return arr
}

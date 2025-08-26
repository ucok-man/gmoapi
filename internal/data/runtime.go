package data

import (
	"fmt"
	"strconv"
)

// Declare a custom Runtime type, yang menggunakan underlying type int32
type Runtime int32

// Implement a MarshalJSON() method pada Runtime type
// Method ini akan membuat output string JSON dengan format "<runtime> mins"
func (r Runtime) MarshalJSON() ([]byte, error) {
	// Buat string runtime dengan format yang diinginkan
	jsonValue := fmt.Sprintf("%d mins", r)

	// Gunakan strconv.Quote() untuk membungkus string dalam tanda kutip
	quotedJSONValue := strconv.Quote(jsonValue)

	// Convert string ke byte slice dan return
	return []byte(quotedJSONValue), nil
}

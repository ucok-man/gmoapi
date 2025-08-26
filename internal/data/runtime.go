package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

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

// Implementasi UnmarshalJSON agar Runtime bisa diparse dari string "<runtime> mins"
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	// Hilangkan tanda kutip (" ") dari JSON string
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Pisahkan string jadi bagian angka dan kata "mins"
	parts := strings.Split(unquotedJSONValue, " ")

	// Validasi format, harus tepat 2 bagian: "<angka> mins"
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	// Ubah angka string jadi int32
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Assign ke pointer Runtime (pakai dereference *)
	*r = Runtime(i)

	return nil
}

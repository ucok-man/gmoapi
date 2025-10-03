// SPDX-FileCopyrightText: The go-mail Authors
//
// SPDX-License-Identifier: MIT

package mail

import (
	"errors"
	"io"
)

// newlineBytes is a byte slice representation of the SingleNewLine constant used for line breaking
// in encoding processes.
var newlineBytes = []byte(SingleNewLine)

// base64LineBreaker handles base64 encoding with the insertion of new lines after a certain number
// of characters.
//
// This struct is used to manage base64 encoding while ensuring that new lines are inserted after
// reaching a specific line length. It satisfies the io.WriteCloser interface.
//
// References:
//   - https://datatracker.ietf.org/doc/html/rfc2045 (Base64 and line length limitations)
type base64LineBreaker struct {
	line [MaxBodyLength]byte
	used int
	out  io.Writer
}

// Write writes data to the base64LineBreaker, ensuring lines do not exceed MaxBodyLength.
//
// This method writes the provided data to the base64LineBreaker. It ensures that the written
// lines do not exceed the MaxBodyLength. If the data exceeds the limit, it handles the
// continuation by splitting the data and writing new lines as necessary.
//
// Parameters:
//   - data: A byte slice containing the data to be written.
//
// Returns:
//   - numBytes: The number of bytes written.
//   - err: An error if one occurred during the write operation.
func (l *base64LineBreaker) Write(data []byte) (numBytes int, err error) {
	if l.out == nil {
		err = errors.New("no io.Writer set for base64LineBreaker")
		return numBytes, err
	}
	if l.used+len(data) < MaxBodyLength {
		copy(l.line[l.used:], data)
		l.used += len(data)
		return len(data), nil
	}

	_, err = l.out.Write(l.line[0:l.used])
	if err != nil {
		return numBytes, err
	}
	excess := MaxBodyLength - l.used
	l.used = 0

	numBytes, err = l.out.Write(data[0:excess])
	if err != nil {
		return numBytes, err
	}

	_, err = l.out.Write(newlineBytes)
	if err != nil {
		return numBytes, err
	}

	var n int
	n, err = l.Write(data[excess:]) // recurse
	numBytes += n
	return numBytes, err
}

// Close finalizes the base64LineBreaker, writing any remaining buffered data and appending a newline.
//
// This method ensures that any remaining data in the buffer is written to the output and appends
// a newline. It is used to finalize the base64LineBreaker and should be called when no more data
// is expected to be written.
//
// Returns:
//   - err: An error if one occurred during the final write operation.
func (l *base64LineBreaker) Close() (err error) {
	if l.used > 0 {
		_, err = l.out.Write(l.line[0:l.used])
		if err != nil {
			return err
		}
		_, err = l.out.Write(newlineBytes)
	}

	return err
}

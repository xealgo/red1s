package resp

import (
	"fmt"

	"github.com/xealgo/red1s"
)

// Error constants
const (
	ErrEmptyPayload = red1s.Error("Empty payload")
	ErrEndOfInput   = red1s.Error("End of input")
)

// Decode decodes a RESP (small subset) encoded message.
func Decode(b []byte) ([]string, error) {
	size := len(b)
	if size == 0 {
		return nil, fmt.Errorf("Error decoding: %w", ErrEmptyPayload)
	}

	tokens := []string{}
	offset := 0

	for i := 0; i < size-1; i++ {
		if b[i] == '\r' && b[i+1] == '\n' {
			if i != 0 {
				tokens = append(tokens, string(b[offset:i]))
			}

			i++
			offset = i + 1
			continue
		}
	}

	return tokens, nil
}

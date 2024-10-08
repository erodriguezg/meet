package base64util

import (
	"encoding/base64"
	"fmt"
)

func Decode(textB64 *string) ([]byte, error) {
	if textB64 == nil {
		return nil, nil
	}
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(*textB64)))
	n, err := base64.StdEncoding.Decode(dst, []byte(*textB64))
	if err != nil {
		return nil, fmt.Errorf("error decoding base64: %w", err)
	}
	return dst[:n], nil
}

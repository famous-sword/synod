package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

func SumHash(reader io.Reader) string {
	h := sha256.New()
	n, err := io.Copy(h, reader)

	fmt.Printf("written: %d, error: %v\n", n, err)

	return hex.EncodeToString(h.Sum(nil))
}

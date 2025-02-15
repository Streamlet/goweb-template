package utility

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	math_rand "math/rand"
)

func RandBytes(length int) []byte {
	b := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err == nil {
		return b
	}
	for i := 0; i < length; i++ {
		b[i] = byte(math_rand.Intn(256))
	}
	return b
}

func RandDigits(length int) string {
	b := RandBytes(length)
	for i := 0; i < len(b); i++ {
		b[i] = '0' + (b[i] % 10)
	}
	return string(b)
}

func RandString(length int) string {
	b := RandBytes((length + 1) / 2)
	s := hex.EncodeToString(b)
	return s[:length]
}

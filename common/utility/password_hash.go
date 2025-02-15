package utility

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
)

func EncryptPassword(password string) (string, string) {
	salt := RandBytes(sha512.Size)
	hash1 := sha512.Sum512([]byte(password))
	hash2 := sha512.Sum512(append(hash1[:], salt...))
	return hex.EncodeToString(salt), hex.EncodeToString(hash2[:])
}

func VerifyPassword(password, salt, hash string) error {
	saltBin, err := hex.DecodeString(salt)
	if err != nil {
		return err
	}
	hash1 := sha512.Sum512([]byte(password))
	hash2 := sha512.Sum512(append(hash1[:], saltBin...))
	if hex.EncodeToString(hash2[:]) != hash {
		return errors.New("password not match")
	}
	return nil
}

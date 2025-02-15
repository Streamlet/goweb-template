package utility

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func ParsePkcs8PrivateKey(privateKeyString string) (*rsa.PrivateKey, error) {
	pemBlock, _ := pem.Decode([]byte(privateKeyString))
	if pemBlock == nil {
		return nil, errors.New("no PEM data is found")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not a RSA private key")
	}
	return key, nil
}

func ParsePkcs8PublicKey(publicKeyString string) (*rsa.PublicKey, error) {
	pemBlock, _ := pem.Decode([]byte(publicKeyString))
	if pemBlock == nil {
		return nil, errors.New("no PEM data is found")
	}
	publicKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not a RSA public key")
	}
	return key, nil
}

func RsaWithSha256Sign(privateKeyString string, content []byte) ([]byte, error) {
	privateKey, err := ParsePkcs8PrivateKey(privateKeyString)
	if err != nil {
		return nil, err
	}

	sha256Hash := sha256.Sum256(content)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, sha256Hash[:])
	return signature, err
}

func RsaWithSha256Verify(publicKeyString string, content []byte, signature []byte) error {
	publicKey, err := ParsePkcs8PublicKey(publicKeyString)
	if err != nil {
		return err
	}

	sha256Hash := sha256.Sum256(content)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, sha256Hash[:], signature)
}

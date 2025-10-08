package transip

import (
	"crypto"
	cryptrand "crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

func loadPrivateKey(path string) (crypto.Signer, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	signer, ok := key.(crypto.Signer)
	if !ok {
		return nil, fmt.Errorf("not a crypto.Signer key")
	}
	return signer, nil
}

func sign(data []byte, privateKeyPath string) (string, error) {
	signer, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		return "", err
	}

	hash := sha512.Sum512(data)
	sig, err := signer.Sign(cryptrand.Reader, hash[:], crypto.SHA512)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(sig)
	return signature, nil
}

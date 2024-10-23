package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/asn1"
	"errors"
	"math/big"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// TODO: implement RSA and ECDSA signing ...

// RSASigner implements the Signer interface for RSA.
type RSASigner struct {
	PrivateKey *rsa.PrivateKey
}

// ECCSigner implements the Signer interface for ECC.
type ECCSigner struct {
	PrivateKey *ecdsa.PrivateKey
}

// Sign generates an RSA signature for the given data using the RSA private key.
func (s *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	// Hash the data using SHA-256 before signing.
	hashed := sha256.Sum256(dataToBeSigned)
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.PrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// Sign generates an ECC signature for the given data using the ECC private key.
func (s *ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	// Hash the data using SHA-256 before signing.
	hashed := sha256.Sum256(dataToBeSigned)

	r, sInt, err := ecdsa.Sign(rand.Reader, s.PrivateKey, hashed[:])
	if err != nil {
		return nil, err
	}

	// Combine r and s into a single byte slice for easier handling.
	signature, err := asn1.Marshal(struct {
		R, S *big.Int
	}{r, sInt})
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// GetSigner Helper function to determine the correct signer based on the key type.
func GetSigner(privateKey interface{}) (Signer, error) {
	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		return &RSASigner{PrivateKey: key}, nil
	case *ecdsa.PrivateKey:
		return &ECCSigner{PrivateKey: key}, nil
	default:
		return nil, errors.New("unsupported private key type")
	}
}

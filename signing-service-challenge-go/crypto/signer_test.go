package crypto_test

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	cryptoLib "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"math/big"
	"testing"
)

// Helper function to verify RSA signatures
func verifyRSASignature(pub *rsa.PublicKey, data, signature []byte) bool {
	hashed := sha256.Sum256(data)
	err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], signature)
	return err == nil
}

// Helper function to verify ECC signatures
func verifyECCSignature(pub *ecdsa.PublicKey, data, signature []byte) bool {
	var esig struct {
		R, S *big.Int
	}
	_, err := asn1.Unmarshal(signature, &esig)
	if err != nil {
		return false
	}

	hashed := sha256.Sum256(data)
	return ecdsa.Verify(pub, hashed[:], esig.R, esig.S)
}

// TestRSA_Signer tests RSA signing functionality.
func TestRSA_Signer(t *testing.T) {

	rsaKey, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}
	data := []byte("Test data for RSA signing")

	signer := &cryptoLib.RSASigner{PrivateKey: rsaKey}

	// Act
	signature, err := signer.Sign(data)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	// Assert
	if !verifyRSASignature(&rsaKey.PublicKey, data, signature) {
		t.Fatalf("RSA signature verification failed")
	}
}

// TestECC_Signer tests ECC signing functionality.
func TestECC_Signer(t *testing.T) {
	eccKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate ECC key: %v", err)
	}
	data := []byte("Test data for ECC signing")

	signer := &cryptoLib.ECCSigner{PrivateKey: eccKey}

	signature, err := signer.Sign(data)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	if !verifyECCSignature(&eccKey.PublicKey, data, signature) {
		t.Fatalf("ECC signature verification failed")
	}
}

// TestGetSigner_RSA tests the GetSigner function for RSA keys.
func TestGetSigner_RSA(t *testing.T) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	signer, err := cryptoLib.GetSigner(rsaKey)
	if err != nil {
		t.Fatalf("Failed to get RSA signer: %v", err)
	}

	if _, ok := signer.(*cryptoLib.RSASigner); !ok {
		t.Fatalf("Expected RSASigner, got %T", signer)
	}
}

// TestGetSigner_ECC tests the GetSigner function for ECC keys.
func TestGetSigner_ECC(t *testing.T) {
	eccKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate ECC key: %v", err)
	}

	signer, err := cryptoLib.GetSigner(eccKey)
	if err != nil {
		t.Fatalf("Failed to get ECC signer: %v", err)
	}

	if _, ok := signer.(*cryptoLib.ECCSigner); !ok {
		t.Fatalf("Expected ECCSigner, got %T", signer)
	}
}

// TestGetSigner_UnsupportedKey tests the GetSigner function for unsupported key types.
func TestGetSigner_UnsupportedKey(t *testing.T) {
	unsupportedKey := &x509.Certificate{}

	signer, err := cryptoLib.GetSigner(unsupportedKey)

	if signer != nil {
		t.Fatalf("Expected nil signer for unsupported key, got %T", signer)
	}
	if err == nil {
		t.Fatalf("Expected error for unsupported key type, got nil")
	}
}

package crypto

import (
	"crypto/elliptic"
	"testing"
)

// TestRSAGenerator tests the RSA key generation logic.
func TestRSAGenerator(t *testing.T) {
	generator := &RSAGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}
	if keyPair == nil || keyPair.Private == nil || keyPair.Public == nil {
		t.Fatalf("Expected non-nil RSA key pair, got nil")
	}

	privateKey := keyPair.Private

	if privateKey.N.BitLen() != 512 {
		t.Fatalf("Expected RSA key with 512 bits, got %d bits", privateKey.N.BitLen())
	}
}

// TestECCGenerator tests the ECC key generation logic.
func TestECCGenerator(t *testing.T) {
	generator := &ECCGenerator{}
	keyPair, err := generator.Generate()

	if err != nil {
		t.Fatalf("Failed to generate ECC key pair: %v", err)
	}
	if keyPair == nil || keyPair.Private == nil || keyPair.Public == nil {
		t.Fatalf("Expected non-nil ECC key pair, got nil")
	}
	privateKey := keyPair.Private

	if privateKey.Curve != elliptic.P384() {
		t.Fatalf("Expected curve P-384, got %v", privateKey.Curve)
	}
}

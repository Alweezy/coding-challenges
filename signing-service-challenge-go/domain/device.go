package domain

// TODO: signature device domain model ...
import (
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

type SignatureDevice struct {
	mu               sync.Mutex
	ID               uuid.UUID
	Label            string
	Algorithm        string
	SignatureCounter int
	LastSignature    []byte
	PrivateKey       interface{}
	PublicKey        interface{}
	Signer           crypto.Signer
}

// BuildSignData generates the secured data string for signing.
// This includes the signature counter, transaction data, and last signature (if any).
func (device *SignatureDevice) BuildSignData(data string) (string, error) {
	device.mu.Lock()
	defer device.mu.Unlock()

	// Use the device ID if no signature has been made yet
	var signature []byte
	if device.SignatureCounter == 0 {
		signature = []byte(device.ID.String())
	} else {
		signature = device.LastSignature
	}

	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	securedData := fmt.Sprintf("%d_%s_%s", device.SignatureCounter, data, encodedSignature)

	return securedData, nil
}

// CommitSignature updates the signature device's state with the newly created signature.
func (device *SignatureDevice) CommitSignature(signature []byte) error {
	device.mu.Lock()
	defer device.mu.Unlock()

	device.SignatureCounter++
	device.LastSignature = signature

	return nil
}

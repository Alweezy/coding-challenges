package service

import (
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/errors"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/infrastructure"
	"github.com/google/uuid"
)

// DeviceService handles operations related to signature devices.
type DeviceService struct {
	deviceRepository infrastructure.DeviceRepository
}

// NewDeviceService creates a new DeviceService.
func NewDeviceService(deviceRepository infrastructure.DeviceRepository) *DeviceService {
	return &DeviceService{deviceRepository: deviceRepository}
}

// CreateSignatureDevice creates and stores a new signature device.
func (s *DeviceService) CreateSignatureDevice(algorithm, label string) (*domain.SignatureDevice, error) {
	deviceID := uuid.New()
	device := &domain.SignatureDevice{
		ID:               deviceID,
		Label:            label,
		Algorithm:        algorithm,
		SignatureCounter: 0,
	}

	// Generate algorithm-based KeyPair
	var err error
	switch algorithm {
	case "RSA":
		generator := &crypto.RSAGenerator{}
		keyPair, err := generator.Generate()
		if err != nil {
			return nil, errors.WrapError(
				err,
				"Failed to generate RSA key pair",
				http.StatusInternalServerError,
			)
		}
		device.PrivateKey = keyPair.Private
		device.PublicKey = keyPair.Public

	case "ECC":
		generator := &crypto.ECCGenerator{}
		keyPair, err := generator.Generate()
		if err != nil {
			return nil, errors.WrapError(
				err,
				"Failed to generate ECC key pair",
				http.StatusInternalServerError,
			)
		}
		device.PrivateKey = keyPair.Private
		device.PublicKey = keyPair.Public

	default:
		return nil, errors.WrapError(
			nil,
			"Unsupported algorithm "+algorithm,
			http.StatusBadRequest,
		)
	}

	// Get signer based on the private key
	device.Signer, err = crypto.GetSigner(device.PrivateKey)
	if err != nil {
		return nil, errors.WrapError(
			err,
			"Failed to get signer for device",
			http.StatusInternalServerError,
		)
	}

	// Save the device in the repository
	err = s.deviceRepository.Save(device.ID.String(), device)
	if err != nil {
		return nil, errors.WrapError(
			err,
			"Failed to save device in repository",
			http.StatusInternalServerError,
		)
	}

	return device, nil
}

// GetDevice retrieves a signature device by ID.
func (s *DeviceService) GetDevice(id string) (*domain.SignatureDevice, bool) {
	device, exists := s.deviceRepository.GetDeviceById(id)
	if !exists {
		return nil, false
	}
	return device, true
}

// ListDevices retrieves all signature devices.
func (s *DeviceService) ListDevices() ([]*domain.SignatureDevice, error) {
	devices, err := s.deviceRepository.GetAllDevices()
	if err != nil {
		return nil, errors.WrapError(
			err,
			"Failed to list devices from repository",
			http.StatusInternalServerError,
		)
	}
	return devices, nil
}

package infrastructure

import "github.com/fiskaly/coding-challenges/signing-service-challenge/domain"

type DeviceRepository interface {
	Save(id string, device *domain.SignatureDevice) error // Now returns an error
	GetDeviceById(id string) (*domain.SignatureDevice, bool)
	UpdateDevice(device *domain.SignatureDevice) error
	GetAllDevices() ([]*domain.SignatureDevice, error)
}

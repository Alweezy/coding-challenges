package mocks

import (
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// MockDeviceRepository is a simple mock implementation of the DeviceRepository interface.
type MockDeviceRepository struct {
	SavedDevices   map[string]*domain.SignatureDevice
	GetDeviceCalls []string
}

// NewMockDeviceRepository creates and returns a new instance of MockDeviceRepository.
func NewMockDeviceRepository() *MockDeviceRepository {
	return &MockDeviceRepository{
		SavedDevices:   make(map[string]*domain.SignatureDevice),
		GetDeviceCalls: []string{},
	}
}

// Save adds a new device to the mock store.
func (m *MockDeviceRepository) Save(id string, device *domain.SignatureDevice) error {
	if _, exists := m.SavedDevices[id]; exists {
		return fmt.Errorf("device with id %s already exists", id)
	}
	m.SavedDevices[id] = device
	return nil
}

// GetDeviceById retrieves a device by its ID.
func (m *MockDeviceRepository) GetDeviceById(id string) (*domain.SignatureDevice, bool) {
	m.GetDeviceCalls = append(m.GetDeviceCalls, id)
	device, exists := m.SavedDevices[id]
	return device, exists
}

// UpdateDevice updates an existing device in the mock store.
func (m *MockDeviceRepository) UpdateDevice(device *domain.SignatureDevice) error {
	if _, exists := m.SavedDevices[device.ID.String()]; !exists {
		return fmt.Errorf("device with id %s not found", device.ID.String())
	}
	m.SavedDevices[device.ID.String()] = device
	return nil
}

// GetAllDevices returns all stored devices.
func (m *MockDeviceRepository) GetAllDevices() ([]*domain.SignatureDevice, error) {
	var devices []*domain.SignatureDevice
	for _, device := range m.SavedDevices {
		devices = append(devices, device)
	}
	return devices, nil
}

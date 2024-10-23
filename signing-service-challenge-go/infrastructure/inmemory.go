package infrastructure

// TODO: in-memory infrastructure ...
import (
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// InMemoryRepository provides thread-safe in-memory storage for signature devices.
type InMemoryRepository struct {
	mu      sync.RWMutex
	devices map[string]*domain.SignatureDevice
}

// NewInMemoryRepository initializes a new InMemoryRepository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		devices: make(map[string]*domain.SignatureDevice),
	}
}

// Save adds a new device to the store in a thread-safe manner.
func (s *InMemoryRepository) Save(id string, device *domain.SignatureDevice) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.devices[id]; exists {
		return fmt.Errorf("device with id %s already exists", id)
	}

	s.devices[id] = device
	return nil
}

// GetDeviceById retrieves a device by its ID in a thread-safe manner.
func (s *InMemoryRepository) GetDeviceById(id string) (*domain.SignatureDevice, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	device, exists := s.devices[id]
	return device, exists
}

// UpdateDevice updates the state of an existing device in the store.
func (s *InMemoryRepository) UpdateDevice(device *domain.SignatureDevice) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.devices[device.ID.String()]; !exists {
		return fmt.Errorf("device with id %s not found", device.ID)
	}
	s.devices[device.ID.String()] = device
	return nil
}

// GetAllDevices returns a slice of all devices in the store.
func (s *InMemoryRepository) GetAllDevices() ([]*domain.SignatureDevice, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var devices []*domain.SignatureDevice
	for _, device := range s.devices {
		devices = append(devices, device)
	}

	return devices, nil
}

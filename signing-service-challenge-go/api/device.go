package api

import (
	"encoding/json"
	"net/http"
)

// CreateSignatureDeviceRequest represents the request to create a signature device.
type CreateSignatureDeviceRequest struct {
	Algorithm string `json:"algorithm"`
	Label     string `json:"label,omitempty"`
}

// CreateSignatureDeviceResponse represents the response after creating a signature device.
type CreateSignatureDeviceResponse struct {
	ID string `json:"id"`
}

// DeviceResponse represents a device's details in the ListDevices response.
type DeviceResponse struct {
	ID               string `json:"id"`
	Label            string `json:"label,omitempty"`
	Algorithm        string `json:"algorithm"`
	SignatureCounter int    `json:"signature_counter"`
}

// ListDevicesResponse represents the response after listing devices.
type ListDevicesResponse struct {
	Devices []DeviceResponse `json:"devices"`
}

// CreateSignatureDevice creates a new signature device.
func (s *Server) CreateSignatureDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	var req CreateSignatureDeviceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid request payload"})
		return
	}

	device, err := s.DeviceService.CreateSignatureDevice(req.Algorithm, req.Label)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{err.Error()})
		return
	}

	response := CreateSignatureDeviceResponse{
		ID: device.ID.String(),
	}

	WriteAPIResponse(w, http.StatusCreated, response)
}

// ListDevices lists all devices.
func (s *Server) ListDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	devices, err := s.DeviceService.ListDevices()
	if err != nil {
		WriteInternalError(w)
		return
	}

	deviceResponses := make([]DeviceResponse, len(devices))
	for i, device := range devices {
		deviceResponses[i] = DeviceResponse{
			ID:               device.ID.String(),
			Label:            device.Label,
			Algorithm:        device.Algorithm,
			SignatureCounter: device.SignatureCounter,
		}
	}

	WriteAPIResponse(w, http.StatusOK, ListDevicesResponse{
		Devices: deviceResponses,
	})
}

// GetDeviceById fetches a specific device by its ID.
func (s *Server) GetDeviceById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	deviceId := r.PathValue("deviceId")

	device, exists := s.DeviceService.GetDevice(deviceId)
	if !exists {
		WriteErrorResponse(w, http.StatusNotFound, []string{"Device not found"})
		return
	}

	response := DeviceResponse{
		ID:               device.ID.String(),
		Label:            device.Label,
		Algorithm:        device.Algorithm,
		SignatureCounter: device.SignatureCounter,
	}

	WriteAPIResponse(w, http.StatusOK, response)
}

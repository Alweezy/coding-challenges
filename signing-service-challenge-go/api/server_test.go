package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/mocks"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func setupServer() *api.Server {
	deviceRepo := mocks.NewMockDeviceRepository()
	deviceService := service.NewDeviceService(deviceRepo)
	transactionService := service.NewTransactionService(deviceRepo)
	return api.NewServer(":8086", deviceRepo, deviceService, transactionService)
}

func setupRouter(s *api.Server) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/api/v0/health", http.HandlerFunc(s.Health))
	mux.Handle("/api/v0/devices", http.HandlerFunc(s.CreateSignatureDevice))
	mux.Handle("/api/v0/devices/list", http.HandlerFunc(s.ListDevices))
	mux.Handle("/api/v0/devices/", http.HandlerFunc(s.GetDeviceById))
	mux.Handle("/api/v0/transactions/", http.HandlerFunc(s.SignTransaction))
	return mux
}

type apiResponse struct {
	Data api.CreateSignatureDeviceResponse `json:"data"`
}

var signResponse struct {
	Data api.SignTransactionResponse `json:"data"`
}

type WrappedListDevicesResponse struct {
	Data api.ListDevicesResponse `json:"data"`
}

func createSignatureDeviceWithServer(t *testing.T, s *api.Server, algorithm, label string, expectedStatus int) string {
	createRequest := api.CreateSignatureDeviceRequest{
		Algorithm: algorithm,
		Label:     label,
	}
	requestBody, err := json.Marshal(createRequest)
	if err != nil {
		t.Fatalf("Error marshalling create signature device request: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Use the router associated with the fresh server instance
	router := setupRouter(s)
	router.ServeHTTP(w, req)

	if w.Code != expectedStatus {
		t.Fatalf("Expected status code %d, got %d", expectedStatus, w.Code)
	}

	var response apiResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	return response.Data.ID
}
func TestCreateSignatureDeviceSuccess(t *testing.T) {
	testCases := []struct {
		name           string
		algorithm      string
		label          string
		expectedStatus int
	}{
		{
			name:           "Valid ECC Device",
			algorithm:      "ECC",
			label:          "Test ECC Device",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Valid RSA Device",
			algorithm:      "RSA",
			label:          "Test RSA Device",
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := setupServer()

			deviceId := createSignatureDeviceWithServer(t, s, tc.algorithm, tc.label, tc.expectedStatus)

			_, exists := s.DeviceRepository.GetDeviceById(deviceId)
			if !exists {
				t.Fatalf("Expected device with ID %s to be stored", deviceId)
			}
		})
	}
}
func TestConcurrentDeviceCreation(t *testing.T) {
	var wg sync.WaitGroup
	s := setupServer()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			label := fmt.Sprintf("Test Device %d", i)
			createSignatureDeviceWithServer(t, s, "RSA", label, http.StatusCreated)
		}(i)
	}

	wg.Wait()
}
func TestCreateSignatureDeviceFailureCases(t *testing.T) {
	testCases := []struct {
		name           string
		algorithm      string
		label          string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Missing Algorithm",
			algorithm:      "",
			label:          "Missing Algorithm Device",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Unsupported algorithm ",
		},
		{
			name:           "Unsupported Algorithm",
			algorithm:      "AES",
			label:          "Unsupported Algorithm Device",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Unsupported algorithm AES",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := setupServer()

			req := api.CreateSignatureDeviceRequest{
				Algorithm: tc.algorithm,
				Label:     tc.label,
			}
			requestBody, err := json.Marshal(req)
			if err != nil {
				t.Fatalf("Error marshalling request: %v", err)
			}

			request := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewBuffer(requestBody))
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router := setupRouter(s)
			router.ServeHTTP(w, request)

			if w.Code != tc.expectedStatus {
				t.Fatalf("Expected status code %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.expectedStatus == http.StatusBadRequest {
				var errResponse map[string][]string
				if err := json.NewDecoder(w.Body).Decode(&errResponse); err != nil {
					t.Fatalf("Error decoding error response: %v", err)
				}

				errors, ok := errResponse["errors"]
				if !ok || len(errors) == 0 {
					t.Fatalf("Expected an 'errors' field with at least one error message")
				}

				if errors[0] != tc.expectedError {
					t.Fatalf("Expected error message '%s', got '%s'", tc.expectedError, errors[0])
				}
			}
		})
	}
}
func TestSignTransaction(t *testing.T) {
	s := setupServer()

	deviceId := createSignatureDeviceWithServer(t, s, "RSA", "Test Device", http.StatusCreated)

	signRequest := api.SignTransactionRequest{
		Data: "Test transaction data",
	}
	signRequestBody, err := json.Marshal(signRequest)
	if err != nil {
		t.Fatalf("Error marshalling sign transaction request: %v", err)
	}

	signReq, err := http.NewRequest(http.MethodPost, "/api/v0/transactions/{deviceId}/sign", bytes.NewBuffer(signRequestBody))
	if err != nil {
		t.Fatalf("Error creating sign transaction request: %v", err)
	}

	signReq.SetPathValue("deviceId", deviceId)
	signReq.Header.Set("Content-Type", "application/json")

	signW := httptest.NewRecorder()

	router := setupRouter(s)
	router.ServeHTTP(signW, signReq)

	if signW.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, signW.Code)
	}

	var signResponse struct {
		Data api.SignTransactionResponse `json:"data"`
	}
	if err := json.NewDecoder(signW.Body).Decode(&signResponse); err != nil {
		t.Fatalf("Error decoding sign transaction response: %v", err)
	}

	if signResponse.Data.SignedData == "" || signResponse.Data.Signature == "" {
		t.Fatalf("Expected non-empty signed data and signature")
	}
}
func TestConcurrentTransactionSigning(t *testing.T) {
	s := setupServer()
	deviceId := createSignatureDeviceWithServer(t, s, "RSA", "Test Device", http.StatusCreated)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			signRequest := api.SignTransactionRequest{
				Data: fmt.Sprintf("Test transaction data %d", i),
			}
			signRequestBody, err := json.Marshal(signRequest)
			if err != nil {
				t.Fatalf("Error marshalling sign transaction request: %v", err)
			}

			signReq, err := http.NewRequest(http.MethodPost, "/api/v0/transactions/{deviceId}/sign", bytes.NewBuffer(signRequestBody))
			if err != nil {
				t.Fatalf("Error creating sign transaction request: %v", err)
			}

			signReq.SetPathValue("deviceId", deviceId)
			signReq.Header.Set("Content-Type", "application/json")

			signW := httptest.NewRecorder()
			router := setupRouter(s)
			router.ServeHTTP(signW, signReq)

			if signW.Code != http.StatusOK {
				t.Fatalf("Expected status code %d, got %d", http.StatusOK, signW.Code)
			}

			var signResponse struct {
				Data api.SignTransactionResponse `json:"data"`
			}
			if err := json.NewDecoder(signW.Body).Decode(&signResponse); err != nil {
				t.Fatalf("Error decoding sign transaction response: %v", err)
			}
			if signResponse.Data.SignedData == "" || signResponse.Data.Signature == "" {
				t.Fatalf("Expected non-empty signed data and signature")
			}

		}(i)
	}

	wg.Wait()
}
func TestGetDeviceById(t *testing.T) {
	s := setupServer()

	deviceId := createSignatureDeviceWithServer(t, s, "ECC", "Test ECC Device", http.StatusCreated)

	t.Run("Valid Device ID", func(t *testing.T) {
		getReq := httptest.NewRequest(http.MethodGet, "/api/v0/devices/{deviceId}", nil)
		getReq.SetPathValue("deviceId", deviceId)
		getW := httptest.NewRecorder()

		router := setupRouter(s)
		router.ServeHTTP(getW, getReq)

		if getW.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, getW.Code)
		}

		var getResponse apiResponse
		if err := json.NewDecoder(getW.Body).Decode(&getResponse); err != nil {
			t.Fatalf("Error decoding get device response: %v", err)
		}
		if getResponse.Data.ID != deviceId {
			t.Fatalf("Expected device ID %s, got %s", deviceId, getResponse.Data.ID)
		}
	})

	t.Run("Invalid Device ID", func(t *testing.T) {
		getReq := httptest.NewRequest(http.MethodGet, "/api/v0/devices/nonexistent", nil)
		getW := httptest.NewRecorder()

		router := setupRouter(s)
		router.ServeHTTP(getW, getReq)

		if getW.Code != http.StatusNotFound {
			t.Fatalf("Expected status code %d, got %d", http.StatusNotFound, getW.Code)
		}
	})
}
func TestListDevices(t *testing.T) {
	s := setupServer()

	t.Run("No Devices", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v0/devices/list", nil)
		w := httptest.NewRecorder()

		router := setupRouter(s)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
		var response api.ListDevicesResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Error decoding list devices response: %v", err)
		}
		if len(response.Devices) != 0 {
			t.Fatalf("Expected no devices, got %d", len(response.Devices))
		}
	})

	t.Run("Multiple Devices", func(t *testing.T) {
		s := setupServer()
		router := setupRouter(s)
		_ = createSignatureDeviceWithServer(t, s, "RSA", "Test RSA Device", http.StatusCreated)
		_ = createSignatureDeviceWithServer(t, s, "ECC", "Test ECC Device", http.StatusCreated)

		req := httptest.NewRequest(http.MethodGet, "/api/v0/devices/list", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var wrappedResponse WrappedListDevicesResponse
		if err := json.NewDecoder(w.Body).Decode(&wrappedResponse); err != nil {
			t.Fatalf("Error decoding list devices response: %v", err)
		}

		if len(wrappedResponse.Data.Devices) != 2 {
			t.Fatalf("Expected 2 devices, got %d", len(wrappedResponse.Data.Devices))
		}
	})
}

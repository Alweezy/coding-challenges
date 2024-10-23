package api

import (
	"encoding/json"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/infrastructure"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
	"net/http"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress      string
	DeviceService      *service.DeviceService
	TransactionService *service.TransactionService
	DeviceRepository   infrastructure.DeviceRepository
}

// NewServer is a factory to instantiate a new Server.
func NewServer(
	listenAddress string,
	deviceRepository infrastructure.DeviceRepository,
	deviceService *service.DeviceService,
	transactionService *service.TransactionService,
) *Server {
	return &Server{
		listenAddress: listenAddress,
		// TODO: add services / further dependencies here ...
		DeviceRepository:   deviceRepository,
		DeviceService:      deviceService,
		TransactionService: transactionService,
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	mux := http.NewServeMux()

	mux.Handle("/api/v0/health", http.HandlerFunc(s.Health))

	// TODO: register further HandlerFuncs here ...
	mux.Handle("/api/v0/devices", http.HandlerFunc(s.CreateSignatureDevice))
	mux.Handle("/api/v0/devices/list", http.HandlerFunc(s.ListDevices))
	mux.Handle("/api/v0/devices/{deviceId}", http.HandlerFunc(s.GetDeviceById))
	mux.Handle("/api/v0/transactions/{deviceId}/sign", http.HandlerFunc(s.SignTransaction))

	return http.ListenAndServe(s.listenAddress, mux)
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
		return
	}

	w.Write(bytes)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
		return
	}

	w.Write(bytes)
}

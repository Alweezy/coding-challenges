package api

import (
	"encoding/json"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/errors"
	"net/http"
)

// SignTransactionRequest represents the request to sign data with a signature device.
type SignTransactionRequest struct {
	Data string `json:"data"`
}

// SignTransactionResponse represents the response after signing the transaction.
type SignTransactionResponse struct {
	SignedData string `json:"signed_data"`
	Signature  string `json:"signature"`
}

// SignTransaction signs data using the specified signature device.
func (s *Server) SignTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	deviceId := r.PathValue("deviceId")

	var req SignTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid request payload"})
		return
	}

	signedData, signature, err := s.TransactionService.SignTransaction(deviceId, req.Data)

	if err != nil {
		appErr := errors.FromError(err)
		WriteErrorResponse(w, appErr.Code, []string{appErr.Message})
		return
	}

	response := SignTransactionResponse{
		SignedData: signedData,
		Signature:  signature,
	}
	WriteAPIResponse(w, http.StatusOK, response)
}

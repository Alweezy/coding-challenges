package service

import (
	"encoding/base64"
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/errors"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/infrastructure"
	"net/http"
)

// TransactionService handles operations related to transactions.
type TransactionService struct {
	deviceRepository infrastructure.DeviceRepository
}

// NewTransactionService creates a new TransactionService.
func NewTransactionService(deviceRepository infrastructure.DeviceRepository) *TransactionService {
	return &TransactionService{deviceRepository: deviceRepository}
}

// SignTransaction signs data using the specified signature device.
func (s *TransactionService) SignTransaction(deviceId string, data string) (string, string, error) {
	device, exists := s.deviceRepository.GetDeviceById(deviceId)
	if !exists {
		return "", "", errors.WrapError(nil,
			fmt.Sprintf(
				"Device with id %s not found", deviceId,
			),
			http.StatusNotFound,
		)
	}

	securedData, err := device.BuildSignData(data)
	if err != nil {
		return "", "", errors.WrapError(
			err,
			"An error occurred while building secured data",
			http.StatusInternalServerError,
		)
	}

	signature, err := device.Signer.Sign([]byte(securedData))
	if err != nil {
		return "", "", errors.WrapError(err,
			"error while signing the data",
			http.StatusInternalServerError,
		)
	}

	err = device.CommitSignature(signature)
	if err != nil {
		return "", "", errors.WrapError(
			err,
			"An error occurred while committing the signature: %w",
			http.StatusInternalServerError,
		)
	}

	err = s.deviceRepository.UpdateDevice(device)
	if err != nil {
		return "", "", errors.WrapError(
			err,
			"An error occurred while updating device in repository: %w",
			http.StatusInternalServerError,
		)
	}

	return securedData, base64.StdEncoding.EncodeToString(signature), nil
}

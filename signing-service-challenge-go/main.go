package main

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/infrastructure"
)

const (
	ListenAddress = ":8600"
	// TODO: add further configuration parameters here ...
)

func main() {
	deviceRepository := infrastructure.NewInMemoryRepository()

	deviceService := service.NewDeviceService(deviceRepository)
	transactionService := service.NewTransactionService(deviceRepository)
	
	server := api.NewServer(ListenAddress, deviceRepository, deviceService, transactionService)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}

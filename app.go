package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	service "dwimc/service"
)

type envs struct {
	databaseUri  string
	databaseName string
	servicePort  string
	secretApiKey string
}

var ctx = context.Background()
var environment *envs

func init() {
	environment = getEnvs()
}

func main() {
	log.Println("Starting service...")

	isShuttingDown := false

	store := service.Store{Context: ctx}
	err := store.Init(environment.databaseUri, environment.databaseName)

	if err != nil {
		log.Fatal(err.Error())
	}

	service := service.Service{
		Store:        &store,
		SecretApiKey: environment.secretApiKey,
	}

	go func() {
		log.Println("Starting service... DONE")

		if err := service.Start(environment.servicePort); err != nil && !isShuttingDown {
			log.Printf("service error: %s\n", err)
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	log.Println("Shutting down service...")

	isShuttingDown = true

	if err := shutdownService(&service); err != nil {
		log.Println("Failed stopping service:", err)
	}

	if err := shutdownStore(&store); err != nil {
		log.Println("Failed stopping store:", err)
	}

	log.Println("Shutting down service... DONE")
}

func getEnvs() *envs {
	databaseUri := os.Getenv("DATABASE_URI")
	databaseName := os.Getenv("DATABASE_NAME")
	servicePort := os.Getenv("PORT")
	secretApiKey := os.Getenv("SECRET_API_KEY")

	if len(databaseUri) == 0 {
		databaseUri = "mongodb://localhost:27017"
	}

	if len(databaseName) == 0 {
		databaseName = "dwimc"
	}

	if len(servicePort) == 0 {
		servicePort = "1337"
	}

	return &envs{
		databaseUri:  databaseUri,
		databaseName: databaseName,
		servicePort:  servicePort,
		secretApiKey: secretApiKey,
	}
}

func shutdownService(service *service.Service) error {
	cctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return service.Stop(cctx)
}

func shutdownStore(store *service.Store) error {
	cctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return store.Close(cctx)
}

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	service "dwimc/internal"
)

func main() {
	err := godotenv.Load()
	if err != nil {
	  log.Println("Warning - error loading .env file")
	}

	isShuttingDown := false

	params := service.ServiceParams{
		DBUri:  getEnv("DATABASE_URI", "mongodb://localhost:27017"),
		DBName: getEnv("DATABASE_NAME", "dwimc"),
		ApiKey: getEnv("SECRET_API_KEY", "please_change_me_api_key"),
		Port:   getEnv("PORT", "1337"),
	}
	service := service.CreateService(params)

	go func() {
		if err := service.Start(); err != nil && !isShuttingDown {
			log.Printf("service error: %s\n", err)
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	isShuttingDown = true

	if err := service.Stop(); err != nil {
		log.Println("Failed stopping service:", err)
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}

	return value
}

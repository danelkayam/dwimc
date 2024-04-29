package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Service struct {
	repository DevicesRepository
	server     *http.Server
}

type ServiceParams struct {
	DBUri  string
	DBName string
	ApiKey string
	Port   string
}

func CreateService(params ServiceParams) Service {
	repository := CreateDevicesRepository(params.DBUri, params.DBName)

	router := gin.Default()
	apiRouter := router.Group("/api")

	CreateDevicesRouter(repository, apiRouter, validateApiKey(params.ApiKey))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", params.Port),
		Handler: router,
	}

	return Service{
		repository: repository,
		server:     server,
	}
}

func (service *Service) Start() error {
	log.Println("Starting service...")
	defer log.Println("Starting service... DONE")

	err := service.repository.init()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	log.Printf("Lifting service on: %v\n", service.server.Addr)

	return service.server.ListenAndServe()
}

func (service *Service) Stop() error {
	log.Println("Shutting down service...")
	defer log.Println("Shutting down service... DONE")

	err1 := service.shutdownService()
	err2 := service.shutdownStore()

	if err2 != nil {
		return err2
	}

	return err1
}

// helper functions

func (service *Service) shutdownService() error {
	cctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return service.server.Shutdown(cctx)
}

func (service *Service) shutdownStore() error {
	cctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return service.repository.close(cctx)
}

// Middlewares

func validateApiKey(secretApiKey string) gin.HandlerFunc {
	if len(secretApiKey) > 0 {
		return func(c *gin.Context) {
			apiKey := c.Request.Header.Get("X-API-Key")

			if apiKey != secretApiKey {
				c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
					Data:  nil,
					Error: &ErrorResponse{Message: "Unauthorized!"},
				})

				return
			}

			c.Next()
		}
	}

	return func(c *gin.Context) { c.Next() }
}

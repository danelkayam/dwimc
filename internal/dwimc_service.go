package dwimc_service

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type DwimcService struct {
	lock   sync.Mutex
	server *http.Server
}

type ServiceParams struct {
	Port         int
	DatabasePath string
	DebugMode    bool
	SecretApiKey string
}

const stop_timeout = 5 * time.Second

func NewDwimcService(params ServiceParams) *DwimcService {

	// TODO - init repositories
	// TODO - init services
	// TODO - setup routers
	router := gin.Default()
	// router.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "hello dwimc",
	// 	})
	// })

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", params.Port),
		Handler: router,
	}

	return &DwimcService{
		server: server,
	}
}

func (s *DwimcService) Start() error {
	log.Info().Msg("Starting service...")
	log.Debug().Msgf("Starting server on: %v", s.server.Addr)

	log.Info().Msg("Starting service... DONE")
	return s.server.ListenAndServe()
}

func (s *DwimcService) Stop() error {
	log.Info().Msg("Stopping service...")
	defer log.Info().Msg("Stopping service... DONE")

	cctx, cancel := context.WithTimeout(context.Background(), stop_timeout)
	defer cancel()

	return s.server.Shutdown(cctx)
}

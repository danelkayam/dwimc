package service

import (
	"context"
	"dwimc/internal/api"
	"dwimc/internal/database"
	"dwimc/internal/repositories"
	"dwimc/internal/services"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const stop_timeout = 5 * time.Second

type APIService struct {
	params APIServiceParams
	client *mongo.Client
	server *http.Server
}

type APIServiceParams struct {
	Port         int
	DatabaseURI  string
	DatabaseName string
	SecretApiKey string
	DebugMode    bool
}

func NewAPIService(params APIServiceParams) APIService {
	return APIService{
		params: params,
	}
}

func (s *APIService) Start() error {
	log.Info().Msg("Starting service...")

	client, err := database.InitializeDatabase(s.params.DatabaseURI)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize the database")
		return err
	}

	s.client = client

	context := context.Background()

	deviceRepo, err := repositories.NewMongodbDeviceRepository(context, client, s.params.DatabaseName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize device repository")
		return err
	}

	locationRepo, err := repositories.NewMongodbLocationRepository(context, client, s.params.DatabaseName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize location repository")
		return err
	}

	deviceService := services.NewDefaultDeviceService(deviceRepo)
	locationService := services.NewDefaultLocationService(locationRepo)

	router := api.InitializeRouters(
		s.params.DebugMode,
		deviceService,
		locationService,
	)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.params.Port),
		Handler: router,
	}

	log.Debug().Msgf("Starting server on: %v", s.server.Addr)
	log.Info().Msg("Starting service... DONE")
	return s.server.ListenAndServe()
}

func (s *APIService) Stop() error {
	log.Info().Msg("Stopping service...")
	defer log.Info().Msg("Stopping service... DONE")

	var err1, err2 error

	if s.server != nil {
		cctx, cancel := context.WithTimeout(context.Background(), stop_timeout)
		defer cancel()

		err1 = s.server.Shutdown(cctx)
	}

	if s.client != nil {
		cctx, cancel := context.WithTimeout(context.Background(), stop_timeout)
		defer cancel()

		err2 = s.client.Disconnect(cctx)
	}

	if err1 != nil && err2 == nil {
		return err1
	}

	if err1 == nil && err2 != nil {
		return err2
	}

	return err1
}

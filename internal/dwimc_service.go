package dwimc_service

import (
	"context"
	// "dwimc/internal/repositories"
	// "dwimc/internal/services"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

type DwimcService struct {
	params ServiceParams
	db     *sqlx.DB
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
	return &DwimcService{params: params}
}

func (s *DwimcService) Start() error {
	log.Info().Msg("Starting service...")

	db, err := s.initDatabase(s.params.DatabasePath)
	if err != nil {
		return fmt.Errorf("failed to init database: %v", err)
	}

	s.db = db

	// userRepo := repositories.NewSQLUserRepository(db)
	// deviceRepo := repositories.NewSQLDeviceRepository(db)
	// locationRepo := repositories.NewSQLLocationRepository(db)

	// userService := services.NewDefaultUserService(userRepo)
	// deviceService := services.NewDefaultDeviceService(deviceRepo)
	// locationService := services.NewDefaultLocationService(locationRepo)

	// TODO - setup routers
	// TODO - init gin routing

	router := gin.Default()
	// TODO - remove this
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello dwimc",
		})
	})

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.params.Port),
		Handler: router,
	}

	log.Debug().Msgf("Starting server on: %v", s.server.Addr)
	log.Info().Msg("Starting service... DONE")
	return s.server.ListenAndServe()
}

func (s *DwimcService) Stop() error {
	log.Info().Msg("Stopping service...")
	defer log.Info().Msg("Stopping service... DONE")

	var err1, err2 error

	if s.server != nil {
		cctx, cancel := context.WithTimeout(context.Background(), stop_timeout)
		defer cancel()

		err1 = s.server.Shutdown(cctx)
	}

	if s.db != nil {
		err2 = s.db.Close()
	}

	if err1 != nil && err2 == nil {
		return err1
	}

	if err1 == nil && err2 != nil {
		return err2
	}

	return err1
}

func (s *DwimcService) initDatabase(databasePath string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping database: %v", err)
	}

	return db, nil
}

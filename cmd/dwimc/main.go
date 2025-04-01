package main

import (
	"fmt"
	dlog "log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	service "dwimc/internal"
	"dwimc/internal/utils"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		dlog.Fatalf("Failed loading required envs: %v\n", err)
		return
	}

	initLogger(config)

	log.Info().Msg("Starting DWIMC app...")

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	isShutingDown := false

	service := service.NewAPIService(service.APIServiceParams{
		Port:                 config.Port,
		DatabaseURI:          config.DatabaseURI,
		DatabaseName:         config.DatabaseName,
		SecretAPIKey:         config.SecretAPIKey,
		DebugMode:            config.DebugMode,
		LocationHistoryLimit: config.LocationHistoryLimit,
	})

	go func() {
		if err := service.Start(); err != nil && !isShutingDown {
			log.Error().Msgf("Failed to start service: %v", err)
			termChan <- syscall.SIGTERM
		}
	}()

	log.Info().Msg("Starting DWIMC app... DONE")

	<-termChan
	isShutingDown = true

	log.Info().Msg("Shutting down DWIMC app...")

	if err := service.Stop(); err != nil {
		log.Error().Msgf("Failed to stop service: %v", err)
	}

	log.Info().Msg("Shutting down DWIMC app... DONE")
}

type Config struct {
	Port                 int    `mapstructure:"PORT" validate:"gte=1,lte=65535"`
	DatabaseURI          string `mapstructure:"DATABASE_URI" validate:"required,nonempty"`
	DatabaseName         string `mapstructure:"DATABASE_NAME" validate:"required,nonempty"`
	DebugMode            bool   `mapstructure:"DEBUG_MODE"`
	LogOutputType        string `mapstructure:"LOG_OUTPUT_TYPE" validate:"oneof=console json"`
	LogLevel             string `mapstructure:"LOG_LEVEL" validate:"oneof=debug info warn error"`
	SecretAPIKey         string `mapstructure:"SECRET_API_KEY" validate:"omitempty,nonempty"`
	LocationHistoryLimit int    `mapstructure:"LOCATION_HISTORY_LIMIT"`
}

func loadConfig() (*Config, error) {
	_ = godotenv.Load()

	viper.AutomaticEnv()

	err := viper.BindEnv("DATABASE_URI")
	if err != nil {
		return nil, fmt.Errorf("failed to bind env: %w", err)
	}

	err = viper.BindEnv("LOCATION_HISTORY_LIMIT")
	if err != nil {
		return nil, fmt.Errorf("failed to bind env: %w", err)
	}

	viper.SetDefault("PORT", 8080)
	viper.SetDefault("DATABASE_NAME", "dwimc")
	viper.SetDefault("DEBUG_MODE", false)
	viper.SetDefault("LOG_OUTPUT_TYPE", "json")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("SECRET_API_KEY", "")
	viper.SetDefault("LOCATION_HISTORY_LIMIT", 0)

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	validate := utils.GetDefaultValidate()
	if err := validate.Struct(cfg); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrors {
				dlog.Printf("Invalid value for ENV variable: %s (failed on %s)", fieldErr.Field(), fieldErr.Tag())
			}
		}
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func initLogger(config *Config) {
	if config.LogOutputType == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	level, err := zerolog.ParseLevel(strings.ToLower(config.LogLevel))
	if err != nil {
		log.Warn().Msgf("Invalid log level: %s, defaulting to INFO", config.LogLevel)
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
}

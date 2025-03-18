package main

import (
	"fmt"
	dlog "log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		dlog.Fatalf("Failed loading required envs: %v\n", err)
		return
	}

	initLogger(config)

	log.Info().Msg("Starting DWIMC app...")
}

type Config struct {
	Port          int    `mapstructure:"PORT" validate:"gte=1,lte=65535"`
	DatabasePath  string `mapstructure:"DATABASE_PATH" validate:"required,filepath,file"`
	DebugMode     bool   `mapstructure:"DEBUG_MODE"`
	LogOutputType string `mapstructure:"LOG_OUTPUT_TYPE" validate:"oneof=console json"`
	LogLevel      string `mapstructure:"LOG_LEVEL" validate:"oneof=debug info warn error"`
	SecretApiKey  string `mapstructure:"SECRET_API_KEY"`
}

func loadConfig() (*Config, error) {
	_ = godotenv.Load()

	viper.AutomaticEnv()

	viper.BindEnv("DATABASE_PATH")

	viper.SetDefault("PORT", 1337)
	viper.SetDefault("DEBUG_MODE", false)
	viper.SetDefault("LOG_OUTPUT_TYPE", "json")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("SECRET_API_KEY", "")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	validate := validator.New()
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

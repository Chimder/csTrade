package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type EnvVars struct {
	Username       string
	Password       string
	SteamID        string
	SharedSecret   string
	IdentitySecret string
	DeviceID       string
	DbUrl          string
	Debug          bool
	Env            string
	LogLevel       string
}

func LoadEnv() *EnvVars {
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg(".env file not found, using system environment variables")
	}
	cfg := &EnvVars{
		Username:       getEnv("USERNAME", ""),
		Password:       getEnv("PASSWORD", ""),
		SteamID:        getEnv("STEAM_ID", ""),
		SharedSecret:   getEnv("SHARED_SECRET", ""),
		IdentitySecret: getEnv("IDENTITY_SECRET", ""),
		DeviceID:       getEnv("DEVICE_ID", ""),
		DbUrl:          getEnv("DB_URL", ""),
		Env:            getEnv("ENV", "dev"),
		LogLevel:       getEnv("LOG_LEVEL", "debug"),
	}

	log.Info().Str("sID", cfg.SteamID).Str("db", cfg.DbUrl).Msg("ENV")
	return cfg
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	log.Info().Str("use default", defaultVal).Str("for key", key).Msg("ENV")
	return defaultVal
}

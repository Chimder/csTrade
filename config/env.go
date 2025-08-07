package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type EnvVars struct {
	Username       string `env:"USERNAME"`
	Password       string `env:"PASSWORD"`
	SteamID        string `env:"STREAM_ID"`
	SharedSecret   string `env:"SHARED_SECRET"`
	IdentitySecret string `env:"IDENTITY_SECRET"`
	DeviceID       string `env:"DEVICE_ID"`
	DBUrl          string `env:"DB_URL"`
	Debug          bool   `env:"DEBUG"`
	ENV            string `env:"ENV"`
	LOG_LEVEL      string `env:"LOG_LEVEL"`
}

func LoadEnv() *EnvVars {
	_ = godotenv.Load()
	cfg := EnvVars{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}

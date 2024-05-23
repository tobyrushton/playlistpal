package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port          string `envconfig:"PORT" default:":8080"`
	SpotifyID     string `envconfig:"SPOTIFY_ID"`
	SpotifySecret string `envconfig:"SPOTIFY_SECRET"`
}

func loadConfig() (*Config, error) {
	godotenv.Load()
	var cfg Config

	err := envconfig.Process("SPOTIFY", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func MustLoadConfig() *Config {
	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}

	return cfg
}

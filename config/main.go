package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		Server: loadServerConfig(),
	}
}

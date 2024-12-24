package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	WebSocket WebSocketConfig
	TurnCredentials TurnCredentialConfig
}

func NewConfig(envFiles []string) *Config {
	err := godotenv.Load(envFiles...)
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		Server:    loadServerConfig(),
		WebSocket: loadWebSocketConfig(),
		TurnCredentials: loadTurnCredentialConfig(),
	}
}

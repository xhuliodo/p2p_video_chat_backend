package config

import (
	"log"
	"os"
	"time"
)

type TurnCredentialConfig struct {
	Secret      string
	ExpireAfter time.Duration
}

func loadTurnCredentialConfig() TurnCredentialConfig {
	secret := os.Getenv("TURN_SERVER_SECRET")
	if secret==""{
		log.Fatal("could not retrieve TURN_SERVER_SECRET env var")
	}
	return TurnCredentialConfig{
		Secret:      secret,
		ExpireAfter: 8 * time.Hour,
	}
}

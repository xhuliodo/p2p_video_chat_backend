package config

import (
	"log"
	"os"
	"time"
)

type TurnCredentialConfig struct {
	CloudflareTurnTokenId string
	CloudflareApiToken    string
	ExpireAfter           time.Duration
}

func loadTurnCredentialConfig() TurnCredentialConfig {
	cloudflareTurnTokenId := os.Getenv("CLOUDFLARE_TURN_TOKEN_ID")
	if cloudflareTurnTokenId == "" {
		log.Fatal("could not retrieve CLOUDFLARE_TURN_TOKEN_ID env var")
	}

	cloudflareApiToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if cloudflareApiToken == "" {
		log.Fatal("could not retrieve TURN_SERVER_SECRET env var")
	}

	return TurnCredentialConfig{
		CloudflareTurnTokenId: cloudflareTurnTokenId,
		CloudflareApiToken:    cloudflareApiToken,
		ExpireAfter:           8 * time.Hour,
	}
}

package config

import (
	"time"
)

type ServerConfig struct {
	Port            string
	GracefulTimeout time.Duration
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
	IdleTimeout     time.Duration
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port:            ":80",
		GracefulTimeout: time.Second * 15,
		WriteTimeout:    time.Second * 15,
		ReadTimeout:     time.Second * 15,
		IdleTimeout:     time.Second * 60,
	}
}

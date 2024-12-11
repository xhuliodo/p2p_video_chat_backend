package config

import (
	"os"
	"time"
)

type ServerConfig struct {
	SSLCert         string
	SSLCertKey      string
	Port            string
	GracefulTimeout time.Duration
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
	IdleTimeout     time.Duration
}

func loadServerConfig() ServerConfig {
	sslCert := os.Getenv("SSL_CERT")
	if sslCert == "" {
		sslCert = "dev_cert.crt"
	}

	sslCertKey := os.Getenv("SSL_CERT_KEY")
	if sslCertKey == "" {
		sslCertKey = "dev_cert.key"
	}
	return ServerConfig{
		SSLCert:         sslCert,
		SSLCertKey:      sslCertKey,
		Port:            ":8080",
		GracefulTimeout: time.Second * 15,
		WriteTimeout:    time.Second * 15,
		ReadTimeout:     time.Second * 15,
		IdleTimeout:     time.Second * 60,
	}
}

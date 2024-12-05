package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/xhuliodo/p2p_video_chat_backend/config"
	"github.com/xhuliodo/p2p_video_chat_backend/server/handler"
)

type Server struct {
	config *config.Config
	http   *http.Server
}

func NewServer(cfg *config.Config) *Server {
	h := handler.NewHandler(cfg)
	r := h.Configure()

	httpServer := &http.Server{
		Addr: cfg.Server.Port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: cfg.Server.WriteTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	return &Server{
		config: cfg,
		http:   httpServer,
	}
}

func (s *Server) Start() error {
	slog.Info("Starting the server on HTTP PORT: " + s.config.Server.Port)
	return s.http.ListenAndServeTLS(s.config.Server.SSLCert, s.config.Server.SSLCertKey)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

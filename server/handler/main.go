package handler

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/xhuliodo/p2p_video_chat_backend/config"
	ws "github.com/xhuliodo/p2p_video_chat_backend/server/websocket"
)

type Handler struct {
	wsConfig   *config.WebSocketConfig
	turnConfig *config.TurnCredentialConfig
	upgrader   websocket.Upgrader
	hub        *ws.Hub
}


func NewHandler(cfg *config.Config) Handler {
	u := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  cfg.WebSocket.ReadBufferSize,
		WriteBufferSize: cfg.WebSocket.WriteBufferSize,
	}

	m := ws.NewHub()

	return Handler{
		wsConfig:   &cfg.WebSocket,
		turnConfig: &cfg.TurnCredentials,
		upgrader:   u,
		hub:        m,
	}
}

// CORS middleware to allow cross-origin requests
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow your frontend domain
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")    // Allow specific methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Allow specific headers

		// Handle preflight request (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) Configure() http.Handler {
	r := mux.NewRouter()

	r.Use(enableCORS)
	r.HandleFunc("/calls/{id}", h.upgradeConnection).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/healthcheck", h.healthcheck).Methods(http.MethodGet)
	r.HandleFunc("/turn/credentials", h.TurnCredentials).Methods(http.MethodGet)

	return r
}

func (h *Handler) WSShutdown() {
	ctx := context.Background()
	h.hub.Shutdown(ctx)
}

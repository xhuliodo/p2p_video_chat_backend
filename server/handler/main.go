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
	config   *config.WebSocketConfig
	upgrader websocket.Upgrader
	hub      *ws.Hub
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
		config:   &cfg.WebSocket,
		upgrader: u,
		hub:      m,
	}
}

func (h *Handler) Configure() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/calls/{id}", h.upgradeConnection).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/healthcheck", h.healthcheck).Methods(http.MethodGet)

	return r
}

func (h *Handler) WSShutdown() {
	ctx := context.Background()
	h.hub.Shutdown(ctx)
}

package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/xhuliodo/p2p_video_chat_backend/config"
	ws "github.com/xhuliodo/p2p_video_chat_backend/server/websocket"
)

type Handler struct {
	upgrader websocket.Upgrader
	hub      *ws.Hub
}

func NewHandler(cfg *config.Config) Handler {
	u := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	m := ws.NewHub()

	return Handler{
		upgrader: u,
		hub:      m,
	}
}

func (h *Handler) Configure() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/calls/{id}", h.upgradeConnection).Methods(http.MethodGet, http.MethodPost)

	return r
}

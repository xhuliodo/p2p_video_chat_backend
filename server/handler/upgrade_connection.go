package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	ws "github.com/xhuliodo/p2p_video_chat_backend/server/websocket"
)

func (h *Handler) upgradeConnection(w http.ResponseWriter, r *http.Request) {
	var callId string
	vars := mux.Vars(r)
	if val, ok := vars["id"]; ok {
		callId = val
	}
	if callId == "" {
		http.Error(w, "roomId is missing from path params", http.StatusBadRequest)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "failed to upgrade connection", http.StatusInternalServerError)
		log.Println("failed to upgrade connection with err: ", err)
		return
	}

	p := ws.NewConnection(h.config, conn, callId)
	log.Println("created participant")
	h.hub.AddConnection(p)
	log.Println("added participant")

	go p.ReadMessages(h.hub)
	go p.WriteMessages(h.hub)
}

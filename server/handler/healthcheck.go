package handler

import "net/http"

func (h *Handler) healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Up and running"))
}

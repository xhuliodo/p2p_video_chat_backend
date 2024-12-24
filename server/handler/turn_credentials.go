package handler

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/xhuliodo/p2p_video_chat_backend/response"
)

func (h *Handler) TurnCredentials(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "userId is a required  query variable", http.StatusBadRequest)
		return
	}

	// Generate timestamp
	timestamp := time.Now().Add(h.turnConfig.ExpireAfter).Unix()
	username := fmt.Sprintf("%d:%s", timestamp, userId)
	// Sign the generated username with our secret
	s := hmac.New(sha1.New, []byte(h.turnConfig.Secret))
	s.Write([]byte(username))
	password := base64.StdEncoding.EncodeToString(s.Sum(nil))

	res := response.TurnCredential{
		Username:  username,
		Password:  password,
		ExpiresAt: timestamp,
	}
	resJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "something went wrong, please retry later", http.StatusInternalServerError)
		return
	}
	w.Write(resJson)
}

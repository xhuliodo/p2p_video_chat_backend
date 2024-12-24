package handler_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/xhuliodo/p2p_video_chat_backend/config"
	"github.com/xhuliodo/p2p_video_chat_backend/response"
	"github.com/xhuliodo/p2p_video_chat_backend/server/handler"
)

var cfg = config.NewConfig([]string{"../../.env"})
var handlers = handler.NewHandler(cfg)
var mux = handlers.Configure()
var server = httptest.NewServer(mux)

func TestTurnCredentials(t *testing.T) {
	userId, _ := uuid.NewV7()
	req, err := http.NewRequest("GET", server.URL+"/turn/credentials?userId="+userId.String(), nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Could not read response body: %v", err)
	}
	
	var turnCred response.TurnCredential
	if err := json.Unmarshal(body, &turnCred); err != nil {
		t.Fatalf("Could not unmarshall turn credentials response with err: %v", err)
	}

	fmt.Printf("turn cred %v", turnCred)
}

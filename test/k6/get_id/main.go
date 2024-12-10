package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gofrs/uuid"
)

type callId struct {
	id   string
	used int
	mu   sync.Mutex // guards n
}

const noMoreThanParticipants = 3

func (c *callId) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.used > noMoreThanParticipants {
		newId, _ := uuid.NewV7()
		c.id = newId.String()

		c.used = 0
	}
	c.used += 1
	w.Write([]byte(c.id))
}

func main() {
	initialId, _ := uuid.NewV7()
	callIdHandler := callId{id: initialId.String()}
	http.Handle("/call/id", &callIdHandler)
	log.Fatal(http.ListenAndServe(":3030", nil))
}

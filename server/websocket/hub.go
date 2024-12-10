package websocket

import (
	"context"
	"fmt"
	"sync"

	"github.com/xhuliodo/p2p_video_chat_backend/domain"
)

type Hub struct {
	connections   Connections
	eventHandlers map[string]domain.EventHandler
	mu            sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		connections:   make(Connections),
		eventHandlers: domain.NewEventHandlers(),
	}
}

func (h *Hub) RemoveConnection(conn *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.connections[conn]; ok {
		conn.ws.Close()
		delete(h.connections, conn)
	}
}

func (h *Hub) AddConnection(conn *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.connections[conn] = true
}

func (h *Hub) Shutdown(ctx context.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.connections {
		c.ws.Close()
		delete(h.connections, c)
	}
}

func (h *Hub) SendToCallParticipantsExcept(callId, userId string, event domain.Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for conn := range h.connections {
		if conn.callId == callId && conn.userId != userId {
			conn.egress <- event
		}
	}
}

func (h *Hub) SendToParticipant(userId string, event domain.Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for conn := range h.connections {
		if conn.userId == userId {
			conn.egress <- event
		}
	}
}

// routeEvent is used kind of like an internal router for events and their coresponsive handler function
func (h *Hub) RouteEvent(event domain.Event, conn *Connection) error {
	// check is the event is supported
	handler, ok := h.eventHandlers[event.Type]
	if !ok {
		return fmt.Errorf("event type: %s is not supported", event.Type)
	}

	return handler(event, conn, h)
}

package websocket

import (
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

func (m *Hub) RemoveConnection(conn *Connection) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.connections[conn]; ok {
		conn.ws.Close()
		delete(m.connections, conn)
	}
}

func (m *Hub) AddConnection(conn *Connection) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[conn] = true
}

func (m *Hub) SendToCallParticipantsExcept(callId, userId string, event domain.Event) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for conn := range m.connections {
		if conn.callId == callId && conn.userId != userId {
			conn.egress <- event
		}
	}
}

func (m *Hub) SendToParticipant(userId string, event domain.Event) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for conn := range m.connections {
		if conn.userId == userId {
			conn.egress <- event
		}
	}
}

// routeEvent is used kind of like an internal router for events and their coresponsive handler function
func (m *Hub) RouteEvent(event domain.Event, conn *Connection) error {
	// check is the event is supported
	handler, ok := m.eventHandlers[event.Type]
	if !ok {
		return fmt.Errorf("event type: %s is not supported", event.Type)
	}

	return handler(event, conn, m)
}

package websocket

import (
	"fmt"
	"sync"
)

type Hub struct {
	participants Participants
	handlers     map[string]EventHandler
	mu           sync.RWMutex
}

func NewHub() *Hub {
	h := map[string]EventHandler{
		EventNewParticipant:  NewParticipantHandler,
		EventParticipantLeft: ParticipantLeftHandler,
		EventOffer:           OfferHandler,
		EventAnswer:          AnswerHandler,
		EventIceCandidate:    IceCandidateHandler,
		EventReconnect:       ReconnectHandler,
	}
	return &Hub{
		participants: make(Participants),
		handlers:     h,
	}
}

func (m *Hub) RemoveParticipant(participant *Participant) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.participants[participant]; ok {
		participant.connection.Close()
		delete(m.participants, participant)
	}
}

func (m *Hub) AddParticipant(participant *Participant) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.participants[participant] = true
}

// routeEvent is used kind of like an internal router for events and their coresponsive handler function
func (m *Hub) routeEvent(event Event, p *Participant) error {
	// check is the event is supported
	handler, ok := m.handlers[event.Type]
	if !ok {
		return fmt.Errorf("event type: %s is not supported", event.Type)
	}

	return handler(event, p, m)
}

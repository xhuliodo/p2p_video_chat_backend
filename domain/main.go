package domain

import "encoding/json"

type Hub interface {
	SendToCallParticipantsExcept(callId, userId string, event Event)
	SendToParticipant(userId string, event Event)
}

type Participant interface {
	GetUserId() string
	GetCallId() string
	SetUserId(userId string)
}

// Event is the Messages sent over the websocket
// Used to differ between different actions
type Event struct {
	// Type is the message type sent
	Type string `json:"type"`
	// Payload is the data Based on the Type
	Payload json.RawMessage `json:"payload"`
}

// EventHandler is a function signature that is used to affect messages on the socket and triggered
// depending on the type
type EventHandler func(event Event, participant Participant, hub Hub) error

func NewEventHandlers() map[string]EventHandler {
	return map[string]EventHandler{
		EventNewParticipant:  HandleEventNewParticipant,
		EventParticipantLeft: HandleEventParticipantLeft,
		EventOffer:           HandleEventOffer,
		EventAnswer:          HandleEventAnswer,
		EventIceCandidate:    HandleEventIceCandidate,
		EventReconnect:       HandleEventReconnect,
	}
}

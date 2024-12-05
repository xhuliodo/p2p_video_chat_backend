package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/xhuliodo/p2p_video_chat_backend/request"
	"github.com/xhuliodo/p2p_video_chat_backend/response"
)

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
type EventHandler func(event Event, p *Participant, hub *Hub) error

const (
	EventNewParticipant  = "new_participant"
	EventParticipantLeft = "participant_left"
	EventOffer           = "offer"
	EventAnswer          = "answer"
	EventIceCandidate    = "ice_candidate"
	EventReconnect       = "reconnect"
)

func NewParticipantHandler(event Event, p *Participant, hub *Hub) error {
	var newParticipantEvent request.NewParticipant
	if err := json.Unmarshal(event.Payload, &newParticipantEvent); err != nil {
		return fmt.Errorf("could not unmarshall event type: %s with payload: %v with err: %s", EventNewParticipant, event.Payload, err)
	}

	p.userId = newParticipantEvent.UserId

	var broadcast response.NewParticipant
	broadcast.ParticipantId = p.userId
	data, err := json.Marshal(broadcast)
	if err != nil {
		return fmt.Errorf("could not marshall event type: %s with payload: %v with err: %s", EventNewParticipant, broadcast, err)
	}

	outgoing := Event{
		Type:    EventNewParticipant,
		Payload: data,
	}
	for part := range hub.participants {
		if part.callId == p.callId && part.userId != p.userId {
			part.egress <- outgoing
		}
	}

	return nil
}

func ParticipantLeftHandler(event Event, p *Participant, hub *Hub) error {
	var broadcast response.NewParticipant
	broadcast.ParticipantId = p.userId
	data, err := json.Marshal(broadcast)
	if err != nil {
		return fmt.Errorf("could not marshall event type: %s with payload: %v with err: %s", EventParticipantLeft, broadcast, err)
	}

	outgoing := Event{
		Type:    EventParticipantLeft,
		Payload: data,
	}
	for part := range hub.participants {
		if part.callId == p.callId && part.userId != p.userId {
			part.egress <- outgoing
		}
	}

	return nil
}

func OfferHandler(event Event, p *Participant, hub *Hub) error {
	var offerEvent request.Offer
	if err := json.Unmarshal(event.Payload, &offerEvent); err != nil {
		return fmt.Errorf("could not unmarshall event type: %s with payload: %v with err: %s", EventOffer, event.Payload, err)
	}

	var broadcast response.Offer
	broadcast.Offer = offerEvent.Offer
	broadcast.From = p.userId
	data, err := json.Marshal(broadcast)
	if err != nil {
		return fmt.Errorf("could not marshall event type: %s with payload: %v with err: %s", EventOffer, broadcast, err)
	}

	outgoing := Event{
		Type:    EventOffer,
		Payload: data,
	}
	for part := range hub.participants {
		if part.userId == offerEvent.To {
			part.egress <- outgoing
		}
	}

	return nil
}

func AnswerHandler(event Event, p *Participant, hub *Hub) error {
	var answerEvent request.Answer
	if err := json.Unmarshal(event.Payload, &answerEvent); err != nil {
		return fmt.Errorf("could not unmarshall event type: %s with payload: %v with err: %s", EventAnswer, event.Payload, err)
	}

	var broadcast response.Answer
	broadcast.Answer = answerEvent.Answer
	broadcast.From = p.userId
	data, err := json.Marshal(broadcast)
	if err != nil {
		return fmt.Errorf("could not marshall event type: %s with payload: %v with err: %s", EventAnswer, broadcast, err)
	}

	outgoing := Event{
		Type:    EventAnswer,
		Payload: data,
	}
	for part := range hub.participants {
		if part.userId == answerEvent.To {
			part.egress <- outgoing
		}
	}

	return nil
}

func IceCandidateHandler(event Event, p *Participant, hub *Hub) error {
	var iceCandidateEvent request.IceCandidate
	if err := json.Unmarshal(event.Payload, &iceCandidateEvent); err != nil {
		return fmt.Errorf("could not unmarshall event type: %s with payload: %v with err: %s", EventIceCandidate, event.Payload, err)
	}

	var broadcast response.IceCandidate
	broadcast.IceCandidate = iceCandidateEvent.IceCandidate
	broadcast.From = p.userId
	data, err := json.Marshal(broadcast)
	if err != nil {
		return fmt.Errorf("could not marshall event type: %s with payload: %v with err: %s", EventIceCandidate, broadcast, err)
	}

	outgoing := Event{
		Type:    EventIceCandidate,
		Payload: data,
	}
	for part := range hub.participants {
		if part.userId == iceCandidateEvent.To {
			part.egress <- outgoing
		}
	}

	return nil
}

func ReconnectHandler(event Event, p *Participant, hub *Hub) error {
	var reconnectEvent request.Reconnect
	if err := json.Unmarshal(event.Payload, &reconnectEvent); err != nil {
		return fmt.Errorf("could not unmarshall event type: %s with payload: %v with err: %s", EventReconnect, event.Payload, err)
	}

	var broadcast response.Reconnect
	broadcast.From = p.userId
	data, err := json.Marshal(broadcast)
	if err != nil {
		return fmt.Errorf("could not marshall event type: %s with payload: %v with err: %s", EventReconnect, broadcast, err)
	}

	outgoing := Event{
		Type:    EventReconnect,
		Payload: data,
	}
	for part := range hub.participants {
		if part.userId == reconnectEvent.To {
			part.egress <- outgoing
			break
		}
	}

	return nil
}

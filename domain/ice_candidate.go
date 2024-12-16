package domain

import (
	"encoding/json"
	"fmt"

	"github.com/xhuliodo/p2p_video_chat_backend/request"
	"github.com/xhuliodo/p2p_video_chat_backend/response"
)

const EventIceCandidate = "ice_candidate"

func HandleEventIceCandidate(event Event, part Participant, hub Hub) error {
	var iceCandidateEvent request.IceCandidate
	if err := json.Unmarshal(event.Payload, &iceCandidateEvent); err != nil {
		return fmt.Errorf("could not unmarshall event type: %s with payload: %s with err: %s", EventIceCandidate, string(event.Payload), err)
	}

	var broadcast response.IceCandidate
	broadcast.IceCandidate = iceCandidateEvent.IceCandidate
	broadcast.From = part.GetUserId()
	data, err := json.Marshal(broadcast)
	if err != nil {
		return fmt.Errorf("could not marshall event type: %s with payload: %v with err: %s", EventIceCandidate, broadcast, err)
	}

	outgoing := Event{
		Type:    EventIceCandidate,
		Payload: data,
	}
	hub.SendToParticipant(iceCandidateEvent.To, outgoing)

	return nil
}

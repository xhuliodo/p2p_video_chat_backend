package domain

import (
	"encoding/json"
	"fmt"

	"github.com/xhuliodo/p2p_video_chat_backend/request"
	"github.com/xhuliodo/p2p_video_chat_backend/response"
)

const EventNewParticipant = "new_participant"

func HandleEventNewParticipant(event Event, part Participant, hub Hub) error {
	var newParticipantEvent request.NewParticipant
	if err := json.Unmarshal(event.Payload, &newParticipantEvent); err != nil {
		return fmt.Errorf("could not unmarshall event type: %s with payload: %s with err: %s", EventNewParticipant, string(event.Payload), err)
	}

	part.SetUserId(newParticipantEvent.UserId)

	var broadcast response.NewParticipant
	broadcast.ParticipantId = newParticipantEvent.UserId
	data, err := json.Marshal(broadcast)
	if err != nil {
		return fmt.Errorf("could not marshall event type: %s with payload: %v with err: %s", EventNewParticipant, broadcast, err)
	}

	outgoing := Event{
		Type:    EventNewParticipant,
		Payload: data,
	}
	hub.SendToCallParticipantsExcept(part.GetCallId(), part.GetUserId(), outgoing)

	return nil
}

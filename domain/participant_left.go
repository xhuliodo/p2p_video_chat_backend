package domain

import (
	"encoding/json"
	"fmt"

	"github.com/xhuliodo/p2p_video_chat_backend/response"
)

const EventParticipantLeft = "participant_left"

func HandleEventParticipantLeft(event Event, part Participant, hub Hub) error {
	var broadcast response.NewParticipant
	broadcast.ParticipantId = part.GetUserId()
	data, err := json.Marshal(broadcast)
	if err != nil {
		return fmt.Errorf("could not marshall event type: %s with payload: %v with err: %s", EventParticipantLeft, broadcast, err)
	}

	outgoing := Event{
		Type:    EventParticipantLeft,
		Payload: data,
	}
	hub.SendToCallParticipantsExcept(part.GetCallId(), part.GetUserId(), outgoing)

	return nil
}

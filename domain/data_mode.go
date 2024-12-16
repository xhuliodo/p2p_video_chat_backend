package domain

import (
	"encoding/json"
	"fmt"

	"github.com/xhuliodo/p2p_video_chat_backend/request"
	"github.com/xhuliodo/p2p_video_chat_backend/response"
)

const EventDataMode = "data_mode"

func HandleEventDataMode(event Event, part Participant, hub Hub) error {
	var dataModeEvent request.DataMode
	if err := json.Unmarshal(event.Payload, &dataModeEvent); err != nil {
		return fmt.Errorf("could not unmarshall event type: %s with payload: %s with err: %s", EventDataMode, string(event.Payload), err)
	}

	var broadcast response.DataMode
	broadcast.IsLowDataMode = dataModeEvent.IsLowDataMode
	data, err := json.Marshal(broadcast)
	if err != nil {
		return fmt.Errorf("could not marshall event type: %s with payload: %v with err: %s", EventDataMode, broadcast, err)
	}

	outgoing := Event{
		Type:    EventDataMode,
		Payload: data,
	}
	hub.SendToCallParticipantsExcept(part.GetCallId(), part.GetUserId(), outgoing)

	return nil
}

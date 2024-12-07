package domain

import (
	"encoding/json"
	"fmt"

	"github.com/xhuliodo/p2p_video_chat_backend/request"
	"github.com/xhuliodo/p2p_video_chat_backend/response"
)

const EventReconnect = "reconnect"

func HandleEventReconnect(event Event, part Participant, hub Hub) error {
	var reconnectEvent request.Reconnect
	if err := json.Unmarshal(event.Payload, &reconnectEvent); err != nil {
		return fmt.Errorf("could not unmarshall event type: %s with payload: %v with err: %s", EventReconnect, event.Payload, err)
	}

	var broadcast response.Reconnect
	broadcast.From = part.GetUserId()
	data, err := json.Marshal(broadcast)
	if err != nil {
		return fmt.Errorf("could not marshall event type: %s with payload: %v with err: %s", EventReconnect, broadcast, err)
	}

	outgoing := Event{
		Type:    EventReconnect,
		Payload: data,
	}
	hub.SendToParticipant(reconnectEvent.To, outgoing)

	return nil
}

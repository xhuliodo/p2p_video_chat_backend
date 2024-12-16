package domain

import (
	"encoding/json"
	"fmt"

	"github.com/xhuliodo/p2p_video_chat_backend/request"
	"github.com/xhuliodo/p2p_video_chat_backend/response"
)

const EventOffer = "offer"

func HandleEventOffer(event Event, part Participant, hub Hub) error {
	var offerEvent request.Offer
	if err := json.Unmarshal(event.Payload, &offerEvent); err != nil {
		return fmt.Errorf("could not unmarshall event type: %s with payload: %s with err: %s", EventOffer, string(event.Payload), err)
	}

	var broadcast response.Offer
	broadcast.Offer = offerEvent.Offer
	broadcast.DataMode = offerEvent.DataMode
	broadcast.From = part.GetUserId()
	data, err := json.Marshal(broadcast)
	if err != nil {
		return fmt.Errorf("could not marshall event type: %s with payload: %v with err: %s", EventOffer, broadcast, err)
	}

	outgoing := Event{
		Type:    EventOffer,
		Payload: data,
	}
	hub.SendToParticipant(offerEvent.To, outgoing)

	return nil
}

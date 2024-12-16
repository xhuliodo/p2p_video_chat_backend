package domain

import (
	"encoding/json"
	"fmt"

	"github.com/xhuliodo/p2p_video_chat_backend/request"
	"github.com/xhuliodo/p2p_video_chat_backend/response"
)

const EventAnswer string = "answer"

func HandleEventAnswer(event Event, part Participant, hub Hub) error {
	var answerEvent request.Answer
	if err := json.Unmarshal(event.Payload, &answerEvent); err != nil {
		return fmt.Errorf("could not unmarshall event type: %s with payload: %s with err: %s", EventAnswer, string(event.Payload), err)
	}

	var broadcast response.Answer
	broadcast.Answer = answerEvent.Answer
	broadcast.From = part.GetUserId()
	data, err := json.Marshal(broadcast)
	if err != nil {
		return fmt.Errorf("could not marshall event type: %s with payload: %v with err: %s", EventAnswer, broadcast, err)
	}

	outgoing := Event{
		Type:    EventAnswer,
		Payload: data,
	}
	hub.SendToParticipant(answerEvent.To, outgoing)

	return nil
}

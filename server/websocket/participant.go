package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Participants map[*Participant]bool

type Participant struct {
	userId     string
	callId     string
	connection *websocket.Conn
	egress     chan Event
}

var (
	// pongWait is how long we will await a pong response from client
	pongWait = 10 * time.Second
	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	pingInterval = (pongWait * 9) / 10
)

func NewParticipant(callId string, conn *websocket.Conn) *Participant {
	return &Participant{
		callId:     callId,
		connection: conn,
		egress:     make(chan Event),
	}
}

func (p *Participant) pongHandler(pongMsg string) error {
	// Current time + Pong Wait time
	log.Println("pong")
	return p.connection.SetReadDeadline(time.Now().Add(pongWait))
}

// handleAbruptClosure simulates the user sending the EventParticipantLeft if not sent already
func handleAbruptClosure(p *Participant, hub *Hub) {
	e := Event{
		Type: EventParticipantLeft,
	}
	hub.routeEvent(e, p)
}

func (p *Participant) ReadMessages(hub *Hub) {
	defer func() {
		handleAbruptClosure(p, hub)
		hub.RemoveParticipant(p)
	}()

	// Configure Wait time for Pong response, use Current time + pongWait
	// This has to be done here to set the first initial timer.
	if err := p.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	// Configure how to handle Pong responses
	p.connection.SetPongHandler(p.pongHandler)

	for {
		_, payload, err := p.connection.ReadMessage()
		if err != nil {
			// log error only if the closing is unexpected
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}

			break // break loop and close conn
		}

		var req Event
		if err := json.Unmarshal(payload, &req); err != nil {
			log.Printf("error marshalling message: %v", err)
			continue
		}

		log.Printf("received %s event type\n", req.Type)

		// route event
		if err := hub.routeEvent(req, p); err != nil {
			log.Println("could not route event with err:", err)
		}
	}
}

func (p *Participant) WriteMessages(hub *Hub) {
	// Create a ticker that triggers a ping at given interval
	ticker := time.NewTicker(pingInterval)

	defer func() {
		ticker.Stop()

		handleAbruptClosure(p, hub)
		hub.RemoveParticipant(p)
	}()
	for {
		select {
		case event, ok := <-p.egress:
			if !ok {
				return
			}
			log.Printf("sending %s event type\n", event.Type)

			if err := p.connection.WriteJSON(event); err != nil {
				log.Println("failed sending event with err:", err)
			}
		case <-ticker.C:
			log.Println("ping")
			if err := p.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("could not send ping event with err:", err)
				return
			}
		}
	}
}

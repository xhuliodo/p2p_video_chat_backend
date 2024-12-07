package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xhuliodo/p2p_video_chat_backend/config"
	"github.com/xhuliodo/p2p_video_chat_backend/domain"
)

type Connections map[*Connection]bool

type Connection struct {
	ws     *websocket.Conn
	config *config.WebSocketConfig
	userId string
	callId string
	egress chan domain.Event
}

func NewConnection(config *config.WebSocketConfig, ws *websocket.Conn, callId string) *Connection {
	return &Connection{
		ws:     ws,
		config: config,
		callId: callId,
		egress: make(chan domain.Event),
	}
}

func (c *Connection) GetUserId() string {
	return c.userId
}

func (c *Connection) GetCallId() string {
	return c.callId
}

func (c *Connection) SetUserId(userId string) {
	c.userId = userId
}

func (c *Connection) pongHandler(pongMsg string) error {
	// Current time + Pong Wait time
	log.Println("pong")
	return c.ws.SetReadDeadline(time.Now().Add(c.config.PongWait))
}

func (c *Connection) ReadMessages(hub *Hub) {
	defer func() {
		handleAbruptClosure(c, hub)
		hub.RemoveConnection(c)
	}()

	// Configure Wait time for Pong response, use Current time + pongWait
	// This has to be done here to set the first initial timer.
	if err := c.ws.SetReadDeadline(time.Now().Add(c.config.PongWait)); err != nil {
		log.Println(err)
		return
	}
	// Configure how to handle Pong responses
	c.ws.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.ws.ReadMessage()
		if err != nil {
			// log error only if the closing is unexpected
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}

			break // break loop and close conn
		}

		var req domain.Event
		if err := json.Unmarshal(payload, &req); err != nil {
			log.Printf("error marshalling message: %v", err)
			continue
		}

		log.Printf("received %s event type\n", req.Type)

		// route event
		if err := hub.RouteEvent(req, c); err != nil {
			log.Println("could not route event with err:", err)
		}
	}
}

func (c *Connection) WriteMessages(hub *Hub) {
	// Create a ticker that triggers a ping at given interval
	ticker := time.NewTicker(c.config.PingInterval)

	defer func() {
		ticker.Stop()

		handleAbruptClosure(c, hub)
		hub.RemoveConnection(c)
	}()
	for {
		select {
		case event, ok := <-c.egress:
			if !ok {
				return
			}
			log.Printf("sending %s event type\n", event.Type)

			if err := c.ws.WriteJSON(event); err != nil {
				log.Println("failed sending event with err:", err)
			}
		case <-ticker.C:
			log.Println("ping")
			if err := c.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("could not send ping event with err:", err)
				return
			}
		}
	}
}

// handleAbruptClosure simulates the user sending the EventParticipantLeft if not sent already
func handleAbruptClosure(p *Connection, hub *Hub) {
	e := domain.Event{
		Type: domain.EventParticipantLeft,
	}
	hub.RouteEvent(e, p)
}

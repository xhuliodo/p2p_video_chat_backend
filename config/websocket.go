package config

import "time"

type WebSocketConfig struct {
	// Based on what kind of messages are being exchanged
	ReadBufferSize  int
	WriteBufferSize int
	// pongWait is how long we will await a pong response from client
	PongWait time.Duration
	// pingInterval has to be less than pongWait, We cant multiply by 0.7 to get 70% of time
	// Because that can make decimals, so instead *7 / 10 to get 70%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	PingInterval time.Duration
}

func loadWebSocketConfig() WebSocketConfig {
	pongWait := time.Second * 30
	return WebSocketConfig{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		PongWait:        pongWait,
		PingInterval:    (pongWait * 7) / 10,
	}
}

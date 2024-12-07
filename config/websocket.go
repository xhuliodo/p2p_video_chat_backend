package config

import "time"

type WebSocketConfig struct {
	// pongWait is how long we will await a pong response from client
	PongWait time.Duration
	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	PingInterval time.Duration
}

func loadWebSocketConfig() WebSocketConfig {
	pongWait := time.Second * 10
	return WebSocketConfig{
		PongWait:     pongWait,
		PingInterval: (pongWait * 9) / 10,
	}
}

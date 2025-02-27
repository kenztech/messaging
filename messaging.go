package messaging

import (
	"net/http"

	"github.com/kenztech/messaging/broker"
	"github.com/kenztech/messaging/persistence"
	"github.com/kenztech/messaging/ws"
)

// System encapsulates the messaging system
type System struct {
	Hub         *ws.Hub
	Store       persistence.Store
	Broker      broker.Broker
	ServeWs     http.HandlerFunc // WebSocket handler
	SendMessage http.HandlerFunc // Message sending handler
}

// Config holds configuration for the messaging system
type Config struct {
	Store  persistence.Store
	Broker broker.Broker
}

// NewSystem initializes a new messaging system
func NewSystem(cfg Config) *System {
	hub := ws.NewHub(cfg.Broker)
	handler := ws.NewHandler(hub, cfg.Store, cfg.Broker)

	// Start the hub in a goroutine
	go hub.Run()

	return &System{
		Hub:         hub,
		Store:       cfg.Store,
		Broker:      cfg.Broker,
		ServeWs:     handler.ServeWs,
		SendMessage: handler.SendMessage,
	}
}

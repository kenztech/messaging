package messaging

import (
	"github.com/go-chi/chi/v5"
	"github.com/kenztech/messaging/broker"
	"github.com/kenztech/messaging/persistence"
	"github.com/kenztech/messaging/ws"
)

// System encapsulates the messaging system
type Messaging struct {
	Hub     *ws.Hub
	Handler *ws.Handler
	Store   persistence.Store
	Broker  broker.Broker
}

// Config holds configuration for the messaging system
type Config struct {
	Store  persistence.Store
	Broker broker.Broker
}

// NewSystem initializes a new messaging system
func NewMessaging(cfg Config) *Messaging {
	hub := ws.NewHub(cfg.Broker)
	handler := ws.NewHandler(hub, cfg.Store, cfg.Broker)

	go hub.Run()

	return &Messaging{
		Hub:     hub,
		Handler: handler,
		Store:   cfg.Store,
		Broker:  cfg.Broker,
	}
}

// RegisterRoutes adds messaging routes to a Chi router
func (s *Messaging) RegisterRoutes(r chi.Router) {
	r.Get("/ws", s.Handler.ServeWs)
	r.Post("/message/{content}", s.Handler.SendMessage)
}

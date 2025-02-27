package ws

import (
	"log"
	"sync"

	"github.com/kenztech/messaging/broker"
	"github.com/kenztech/messaging/models"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan models.Message
	register   chan *Client
	unregister chan *Client
	broker     broker.Broker
	mu         sync.RWMutex
}

func NewHub(broker broker.Broker) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan models.Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broker:     broker,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.broker.TrackUser(client.userID, client.groupIDs)
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				close(client.send)
				delete(h.clients, client)
				h.broker.UntrackUser(client.userID, client.groupIDs)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			data, err := message.Marshal()
			if err != nil {
				log.Printf("Failed to marshal message: %v", err)
				continue
			}
			if message.TargetID != "" {
				h.broker.Publish("msg:user:"+message.TargetID, data)
			} else if len(message.GroupIDs) > 0 {
				for _, groupID := range message.GroupIDs {
					members, _ := h.broker.GetGroupMembers(groupID)
					for _, userID := range members {
						h.broker.Publish("msg:user:"+userID, data)
					}
				}
			}
		}
	}
}

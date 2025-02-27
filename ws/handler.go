package ws

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/kenztech/messaging/broker"
	"github.com/kenztech/messaging/models"
	"github.com/kenztech/messaging/persistence"
)

type Handler struct {
	hub      *Hub
	store    persistence.Store
	broker   broker.Broker
	upgrader websocket.Upgrader
}

func NewHandler(hub *Hub, store persistence.Store, broker broker.Broker) *Handler {
	return &Handler{
		hub:    hub,
		store:  store,
		broker: broker,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *Handler) ServeWs(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userId query parameter required", http.StatusBadRequest)
		return
	}
	groupIDs := r.URL.Query()["groupId"]

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	client := NewClient(h.hub, conn, userID, groupIDs, h.broker)
	h.hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	senderID := r.URL.Query().Get("senderId")
	if senderID == "" {
		http.Error(w, "senderId query parameter required", http.StatusBadRequest)
		return
	}

	message := models.NewMessage(
		time.Now().String(),
		senderID,
		chi.URLParam(r, "content"),
		r.URL.Query().Get("targetId"),
		r.URL.Query()["groupId"],
		time.Now().Unix(),
	)

	if err := h.store.SaveMessage(message); err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	h.hub.broadcast <- message
	w.WriteHeader(http.StatusAccepted)
}

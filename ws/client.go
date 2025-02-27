package ws

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kenztech/messaging/broker"
	"github.com/kenztech/messaging/models"
)

type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	send      chan models.Message
	userID    string
	groupIDs  []string
	broker    broker.Broker
	closeChan chan struct{}
}

func NewClient(hub *Hub, conn *websocket.Conn, userID string, groupIDs []string, broker broker.Broker) *Client {
	return &Client{
		hub:       hub,
		conn:      conn,
		send:      make(chan models.Message, 256),
		userID:    userID,
		groupIDs:  groupIDs,
		broker:    broker,
		closeChan: make(chan struct{}),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		close(c.closeChan)
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	msgChan, cleanup, err := c.broker.Subscribe("msg:user:" + c.userID)
	if err != nil {
		log.Printf("Failed to subscribe for user %s: %v", c.userID, err)
		return
	}
	defer cleanup()

	go func() {
		for {
			select {
			case msg, ok := <-msgChan:
				if !ok {
					return
				}
				var m models.Message
				if err := m.Unmarshal(msg); err == nil {
					c.send <- m
				}
			case <-c.closeChan:
				return
			}
		}
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read error for user %s: %v", c.userID, err)
			}
			break
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(message); err != nil {
				log.Printf("write error for user %s: %v", c.userID, err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

package models

import (
	"encoding/json"
	"time"
)

// Message represents a chat message
type Message struct {
	ID        string   `json:"id"`
	SenderID  string   `json:"senderId"`
	Content   string   `json:"content"`
	TargetID  string   `json:"targetId,omitempty"` // Optional: specific user
	GroupIDs  []string `json:"groupIds,omitempty"` // Optional: group(s)
	Timestamp int64    `json:"timestamp"`
}

// NewMessage creates a new message instance
func NewMessage(id, senderID, content, targetID string, groupIDs []string, timestamp int64) Message {
	if id == "" {
		id = time.Now().String()
	}
	return Message{
		ID:        id,
		SenderID:  senderID,
		Content:   content,
		TargetID:  targetID,
		GroupIDs:  groupIDs,
		Timestamp: timestamp,
	}
}

// Marshal serializes the Message to JSON
func (m Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshal deserializes JSON into a Message
func (m *Message) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}

package persistence

import "github.com/kenztech/messaging/models"

// Store defines the interface for persistent storage
type Store interface {
	SaveMessage(m models.Message) error
	GetMessages(userID string, groupIDs []string) ([]models.Message, error)
}

package broker

// Broker defines the interface for message brokering
type Broker interface {
	Publish(channel string, message []byte) error
	Subscribe(channel string) (chan []byte, func(), error)
	TrackUser(userID string, groupIDs []string) error
	UntrackUser(userID string, groupIDs []string) error
	GetGroupMembers(groupID string) ([]string, error)
}

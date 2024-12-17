package api




import (
	"time"
)

// User represents the details of a registered user
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Photo    string `json:"photo"`
}

// Conversation represents the details of a conversation between users
type Conversation struct {
	ID         int64     `json:"id"`
	LastConvo  time.Time `json:"lastconvo"`
	Participant User     `json:"participant"`
}

// Message represents the details of a message sent between users or in a group
type Message struct {
	ID            int64          `json:"id"`
	Timestamp     time.Time      `json:"timestamp"`
	Content       string         `json:"content"`
	SenderUsername User          `json:"senderUsername"`
	Status        string         `json:"status"`
	Conversation  *Conversation  `json:"conversation,omitempty"` // Foreign key to Conversation (if applicable)
	Group         *Group         `json:"group,omitempty"`        // Foreign key to Group (if applicable)
}
type Comment struct {
	ID      int64    `json:"id"`
	Content string   `json:"content"`
	Message *Message `json:"message"` // Foreign key to the associated message
}

// Group represents the details and operations available for managing a group
type Group struct {
	ID       int64   `json:"id"`
	Photo    string  `json:"photo"`
	Name     string  `json:"name"`
	Members  []User  `json:"members"`
	MaxItems int     `json:"maxItems"`
}

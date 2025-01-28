package database

import (
	"database/sql"
	"time"
)

type User struct {
	ID       string         `json:"id"`
	Username string         `json:"username"`
	Photo    sql.NullString `json:"photo"`
}

type Conversation struct {
	ID        int            `json:"id"`
	LastConvo time.Time      `json:"last_convo"`
	IsGroup   bool           `json:"is_group"`
	Photo     sql.NullString `json:"photo"`
	Name      string         `json:"name"`
}

type Convmember struct {
	ID             int `json:"id"`
	ConversationID int `json:"conversation_id"`
	UserID         int `json:"user_id"`
}

type Message struct {
	ID             int       `json:"id"`
	Datetime       time.Time `json:"datetime"`
	Content        string    `json:"content"`
	Sender         int       `json:"sender"`
	ConversationID int       `json:"conversation_id"`
	Status         string    `json:"status"`
}

type ConversationInfo struct {
	OtherUserName string    `json:"other_user_name"` // Name of the other user in the conversation
	LastConvo     time.Time `json:"last_convo"`      // Timestamp of the last conversation
}

// type Participant struct {
//     ID       string `json:"id"`
//     Username string `json:"username"`
//     Photo    string `json:"photo"`
// }

type MessageWithSender struct {
	ID             int            `json:"id"`
	Datetime       time.Time      `json:"datetime"`
	Content        string         `json:"content"`
	SenderID       string         `json:"sender_id"`
	SenderUsername string         `json:"sender_username"`
	SenderPhoto    sql.NullString `json:"sender_photo"` // Use sql.NullString for nullable fields
}

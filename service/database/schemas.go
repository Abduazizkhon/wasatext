package database

import (
	"database/sql"
	"time"
)

type User struct {
	ID       int    		`json:"id"`
	Username string 		`json:"username"`
	Photo    sql.NullString `json:"photo"`
	Token	 string			`json:"token"`

}

type UserToken struct {
	ID int `json:"id"`
	Token string `json:"token"`
	User_id int `json:"user_id"`
}

type Conversation struct {
	ID          int       `json:"id"`
	LastConvo   time.Time `json:"last_convo"`
	IsGroup     bool      `json:"is_group"`
	Photo       sql.NullString    `json:"photo"`
	Name		string		`json:"name"`
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



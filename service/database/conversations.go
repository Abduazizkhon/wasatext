package database

import (
	"database/sql"

)

// GetMyConversations_db retrieves all conversations for a specific user based on their UUID (id).
// func (db *appdbimpl) GetMyConversations_db(userID string) (conversations []Conversation, err error) {
// 	query := `
// 		SELECT 
// 			conversations.id, conversations.lastconvo, conversations.is_group, conversations.photo, conversations.name
// 		FROM 
// 			convmembers
// 		JOIN 
// 			conversations ON convmembers.conversation_id = conversations.id
// 		WHERE 
// 			convmembers.user_id = ?;
// 	`

// 	convos, err := db.c.Query(query, userID)
// 	if err != nil {
// 		return
// 	}
// 	defer convos.Close() // Ensure the result set is closed

// 	for convos.Next() {
// 		var convo Conversation
// 		err = convos.Scan(&convo.ID, &convo.LastConvo, &convo.IsGroup, &convo.Photo, &convo.Name)
// 		if err != nil {
// 			return
// 		}
// 		conversations = append(conversations, convo)
// 	}

// 	// Check for errors during iteration
// 	if err = convos.Err(); err != nil {
// 		return
// 	}

// 	return conversations, nil
// }

// CreateConversation_db creates a new conversation and returns its details.


// AddUsersToConversation adds a user to a conversation.


// GetConversationById retrieves the details of a specific conversation by its ID.
func (db *appdbimpl) GetConversationById(conversationID int) (conversation Conversation, err error) {
	query := `
		SELECT id, lastconvo, is_group, photo, name
		FROM conversations
		WHERE id = ?;
	`

	err = db.c.QueryRow(query, conversationID).Scan(
		&conversation.ID,
		&conversation.LastConvo,
		&conversation.IsGroup,
		&conversation.Photo,
		&conversation.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Conversation{}, nil // Return an empty conversation if no rows are found
		}
		return Conversation{}, err
	}

	return conversation, nil
}

// -------Messages-----



func (db *appdbimpl) CreateConversation_db(isGroup bool, name string, photo string) (conversation Conversation, err error) {
	query := `
		INSERT INTO conversations (lastconvo, is_group, name, photo)
		VALUES (current_timestamp, ?, ?, ?)
		RETURNING id, lastconvo, is_group, name, photo;
	`
	err = db.c.QueryRow(query, isGroup, name, photo).Scan(
		&conversation.ID,
		&conversation.LastConvo,
		&conversation.IsGroup,
		&conversation.Name,
		&conversation.Photo,
	)
	return
}
func (db *appdbimpl) AddUsersToConversation(userID string, conversationID int) (err error) {
	query := `
		INSERT INTO convmembers (user_id, conversation_id)
		VALUES (?, ?);
	`
	_, err = db.c.Exec(query, userID, conversationID)
	return
}
// func (db *appdbimpl) SendMessage(conversationID int, senderID string, content string) error {
// 	query := `
// 		INSERT INTO messages (conversation_id, sender, content, datetime, status)
// 		VALUES (?, ?, ?, CURRENT_TIMESTAMP, 'sent');
// 	`
// 	_, err := db.c.Exec(query, conversationID, senderID, content)
// 	return err
// }

// Check if a conversation already exists between two users
func (db *appdbimpl) ConversationExists(senderID string, recipientID string) (bool, error) {
    query := `
        SELECT COUNT(*) 
        FROM convmembers cm1
        JOIN convmembers cm2 ON cm1.conversation_id = cm2.conversation_id
        WHERE cm1.user_id = ? AND cm2.user_id = ?;
    `
    var count int
    err := db.c.QueryRow(query, senderID, recipientID).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

// GetMyConversations_db retrieves all conversations for a specific user.
func (db *appdbimpl) GetMyConversations_db(userID string) ([]ConversationInfo, error) {
    query := `
        SELECT 
            u.name AS other_user_name,
            c.lastconvo
        FROM 
            convmembers cm1
        JOIN 
            convmembers cm2 ON cm1.conversation_id = cm2.conversation_id
        JOIN 
            users u ON cm2.user_id = u.id
        JOIN 
            conversations c ON cm1.conversation_id = c.id
        WHERE 
            cm1.user_id = ? AND cm2.user_id != ?;
    `
    rows, err := db.c.Query(query, userID, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var conversations []ConversationInfo
    for rows.Next() {
        var convo ConversationInfo
        err := rows.Scan(&convo.OtherUserName, &convo.LastConvo)
        if err != nil {
            return nil, err
        }
        conversations = append(conversations, convo)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return conversations, nil
}

// SendMessage inserts a new message into the database.
func (db *appdbimpl) SendMessage(conversationID int, senderID string, content string) error {
    query := `
        INSERT INTO messages (conversation_id, sender, content, datetime, status)
        VALUES (?, ?, ?, CURRENT_TIMESTAMP, 'sent');
    `
    _, err := db.c.Exec(query, conversationID, senderID, content)
    return err
}

func (db *appdbimpl) IsUserInConversation(userID string, conversationID int) (bool, error) {
    query := `
        SELECT COUNT(*) 
        FROM convmembers 
        WHERE user_id = ? AND conversation_id = ?;
    `
    var count int
    err := db.c.QueryRow(query, userID, conversationID).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}


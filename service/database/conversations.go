package database

import (
	"database/sql"
	"errors"
	"io"
	"mime/multipart" 
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	var photoURL sql.NullString
	if photo != "" {
		photoURL = sql.NullString{String: photo, Valid: true} // Set the photo URL correctly
	} else {
		photoURL = sql.NullString{Valid: false} // If there's no photo, set it as NULL
	}
	query := `
		INSERT INTO conversations (lastconvo, is_group, name, photo)
		VALUES (current_timestamp, ?, ?, ?)
		RETURNING id, lastconvo, is_group, name, photo;
	`
	err = db.c.QueryRow(query, isGroup, name, photoURL).Scan(
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
        JOIN conversations c ON cm1.conversation_id = c.id
        WHERE cm1.user_id = ? 
        AND cm2.user_id = ? 
        AND c.is_group = FALSE;
    `

	var count int
	err := db.c.QueryRow(query, senderID, recipientID).Scan(&count)
	if err != nil {
		return false, err // Returning the error, so it gets logged in the caller function
	}

	return count > 0, nil
}

// GetMyConversations_db retrieves all conversations for a specific user.
// GetMyConversations_db retrieves all conversations for a specific user.
// GetMyConversations_db retrieves all conversations for a specific user.
// GetMyConversations_db retrieves all conversations for a specific user.
// GetMyConversations_db retrieves all conversations for a specific user.
func (db *appdbimpl) GetMyConversations_db(userID string) ([]Conversation, error) {
    query := `
        SELECT 
            c.id, 
            c.lastconvo, 
            c.is_group, 
            c.photo,
            CASE 
                WHEN c.is_group = TRUE THEN c.name 
                ELSE (SELECT u.name FROM users u 
                      JOIN convmembers cm ON u.id = cm.user_id 
                      WHERE cm.conversation_id = c.id AND u.id != ? LIMIT 1)
            END AS name,
            CASE
                WHEN c.is_group = FALSE THEN (SELECT u.photo FROM users u 
                                              JOIN convmembers cm ON u.id = cm.user_id 
                                              WHERE cm.conversation_id = c.id AND u.id != ? LIMIT 1)
                ELSE NULL
            END AS user_photo
        FROM 
            conversations c
        JOIN 
            convmembers cm ON c.id = cm.conversation_id
        WHERE 
            cm.user_id = ?;
    `

    rows, err := db.c.Query(query, userID, userID, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var conversations []Conversation
    for rows.Next() {
        var convo Conversation
        var userPhoto sql.NullString
        err := rows.Scan(&convo.ID, &convo.LastConvo, &convo.IsGroup, &convo.Photo, &convo.Name, &userPhoto)
        if err != nil {
            return nil, err
        }

        // Correct the photo URL by avoiding double '/uploads/'
        if userPhoto.Valid && userPhoto.String != "" {
            // Only prepend '/uploads/' if it's not already there
            if !strings.HasPrefix(userPhoto.String, "/uploads/") {
                convo.Photo.String = "http://localhost:3000/uploads/" + userPhoto.String
            } else {
                convo.Photo.String = "http://localhost:3000" + userPhoto.String
            }
        } else if convo.Photo.Valid && convo.Photo.String != "" {
            // Handle the conversation's own photo URL
            if !strings.HasPrefix(convo.Photo.String, "/uploads/") {
                convo.Photo.String = "http://localhost:3000/uploads/" + convo.Photo.String
            } else {
                convo.Photo.String = "http://localhost:3000" + convo.Photo.String
            }
        } else {
            // If no valid photo, use the default profile picture
            convo.Photo.String = "http://localhost:3000/default-profile.png"
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
func (db *appdbimpl) SendMessageFull(conversationID int, senderID string, content string) error {
	query := `
        INSERT INTO messages (conversation_id, sender, content, datetime, status)
        VALUES (?, ?, ?, CURRENT_TIMESTAMP, 'sent');
    `
	_, err := db.c.Exec(query, conversationID, senderID, content)
	return err
}

func (db *appdbimpl) GetMessagesByConversationId(conversationID int) ([]MessageWithSender, error) {
	query := `
        SELECT 
            m.id, 
            m.datetime, 
            m.content, 
            u.id AS sender_id, 
            u.name AS sender_username, 
            u.photo AS sender_photo
        FROM 
            messages m
        JOIN 
            users u ON m.sender = u.id
        WHERE 
            m.conversation_id = ?
        ORDER BY 
            m.datetime ASC;
    `
	rows, err := db.c.Query(query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []MessageWithSender
	for rows.Next() {
		var msg MessageWithSender
		err := rows.Scan(
			&msg.ID,
			&msg.Datetime,
			&msg.Content,
			&msg.SenderID,
			&msg.SenderUsername,
			&msg.SenderPhoto, // Now using sql.NullString
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (db *appdbimpl) IsMessageOwner(userID string, messageID int) (bool, error) {
	query := `
        SELECT COUNT(*) 
        FROM messages 
        WHERE id = ? AND sender = ?;
    `
	var count int
	err := db.c.QueryRow(query, messageID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (db *appdbimpl) DeleteMessage(messageID int) error {
    query := `DELETE FROM messages WHERE id = ?;`
    _, err := db.c.Exec(query, messageID)
    return err
}

func (db *appdbimpl) GetMessageContent(messageID int) (string, error) {
	query := `
        SELECT content 
        FROM messages 
        WHERE id = ?;
    `
	var content string
	err := db.c.QueryRow(query, messageID).Scan(&content)
	if err != nil {
		return "", err
	}
	return content, nil
}

func (db *appdbimpl) ForwardMessage(targetConversationID int, senderID string, content string) error {
	query := `
        INSERT INTO messages (conversation_id, sender, content, datetime, status)
        VALUES (?, ?, ?, CURRENT_TIMESTAMP, 'forwarded');
    `
	_, err := db.c.Exec(query, targetConversationID, senderID, content)
	return err
}

// ✅ Remove a user from a group
func (db *appdbimpl) RemoveUserFromGroup(userID string, groupID int) error {
	query := `DELETE FROM convmembers WHERE user_id = ? AND conversation_id = ?;`
	_, err := db.c.Exec(query, userID, groupID)
	return err
}

// ✅ Get the count of remaining members in a group
func (db *appdbimpl) GetGroupMemberCount(groupID int) (int, error) {
	query := `SELECT COUNT(*) FROM convmembers WHERE conversation_id = ?;`
	var count int
	err := db.c.QueryRow(query, groupID).Scan(&count)
	return count, err
}

// ✅ Delete a group if empty
func (db *appdbimpl) DeleteGroup(groupID int) error {
	query := `DELETE FROM conversations WHERE id = ?;`
	_, err := db.c.Exec(query, groupID)
	return err
}

// ✅ Check if a conversation is a group
func (db *appdbimpl) IsConversationGroup(conversationID int) (bool, error) {
	query := `SELECT is_group FROM conversations WHERE id = ?;`
	var isGroup bool
	err := db.c.QueryRow(query, conversationID).Scan(&isGroup)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // No conversation found
		}
		return false, err
	}
	return isGroup, nil
}

// ✅ Check if a group with the given name already exists
func (db *appdbimpl) GroupNameExists(name string) (bool, error) {
	query := `SELECT COUNT(*) FROM conversations WHERE is_group = TRUE AND name = ?;`
	var count int
	err := db.c.QueryRow(query, name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ✅ Update the name of a group conversation
func (db *appdbimpl) UpdateGroupName(groupID int, newName string) error {
	query := `UPDATE conversations SET name = ? WHERE id = ? AND is_group = TRUE;`
	_, err := db.c.Exec(query, newName, groupID)
	return err
}

// ✅ Update the group photo
func (db *appdbimpl) UpdateGroupPhoto(groupID int, photoPath string) error {
	query := `UPDATE conversations SET photo = ? WHERE id = ? AND is_group = TRUE;`
	_, err := db.c.Exec(query, photoPath, groupID)
	return err
}

// Check if a message exists
func (db *appdbimpl) DoesMessageExist(messageID int) (bool, error) {
    query := `SELECT COUNT(*) FROM messages WHERE id = ?;`
    var count int
    err := db.c.QueryRow(query, messageID).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

// Add a comment to a message
func (db *appdbimpl) CommentOnMessage(messageID int, userID string, contentType string, content string) error {
	query := `
        INSERT INTO message_comments (message_id, user_id, content_type, content, timestamp)
        VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP);
    `
	_, err := db.c.Exec(query, messageID, userID, contentType, content)
	return err
}

func (db *appdbimpl) DoesConversationExist(conversationID int) (bool, error) {
    query := `SELECT COUNT(*) FROM conversations WHERE id = ?;`
    var count int
    err := db.c.QueryRow(query, conversationID).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func (db *appdbimpl) ConvertCommentsToMessages(messageID int, conversationID int) error {
    querySelect := `
        SELECT user_id, content_type, content FROM message_comments WHERE message_id = ?;
    `
    
    // Fetch comments first
    rows, err := db.c.Query(querySelect, messageID)
    if err != nil {
        return err
    }
    
    // Store comments in memory before closing rows
    var comments []struct {
        UserID      string
        ContentType string
        Content     string
    }

    for rows.Next() {
        var comment struct {
            UserID      string
            ContentType string
            Content     string
        }
        if err := rows.Scan(&comment.UserID, &comment.ContentType, &comment.Content); err != nil {
            rows.Close() // Ensure we close before returning an error
            return err
        }
        comments = append(comments, comment)
    }
    
    rows.Close() // ✅ Properly close before inserting new messages

    // Insert each comment as a new message
    for _, comment := range comments {
        queryInsert := `
            INSERT INTO messages (conversation_id, sender, content, datetime, status)
            VALUES (?, ?, ?, CURRENT_TIMESTAMP, 'comment-converted');
        `
        _, err = db.c.Exec(queryInsert, conversationID, comment.UserID, comment.Content)
        if err != nil {
            return err
        }
    }

    // Finally, delete comments from message_comments
    queryDelete := `DELETE FROM message_comments WHERE message_id = ?;`
    _, err = db.c.Exec(queryDelete, messageID)
    
    return err
}

// ✅ Check if a user is the owner of a comment
func (db *appdbimpl) IsCommentOwner(userID string, commentID int) (bool, error) {
	query := `SELECT COUNT(*) FROM message_comments WHERE id = ? AND user_id = ?;`
	var count int
	err := db.c.QueryRow(query, commentID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ✅ Delete a comment
func (db *appdbimpl) DeleteComment(commentID int) error {
	query := `DELETE FROM message_comments WHERE id = ?;`
	_, err := db.c.Exec(query, commentID)
	return err
}



func (db *appdbimpl) SendMessageWithType(conversationID int, senderID string, content string, contentType string) error {
    query := `
        INSERT INTO messages (conversation_id, sender, content, content_type, datetime, status)
        VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, 'sent');
    `
    _, err := db.c.Exec(query, conversationID, senderID, content, contentType)
    return err
}

func (db *appdbimpl) SendMessageWithMedia(conversationID int, senderID string, contentType string, content string) error {
	query := `
        INSERT INTO messages (conversation_id, sender, content, content_type, datetime, status)
        VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, 'sent');
    `
	_, err := db.c.Exec(query, conversationID, senderID, content, contentType)
	return err
}

// SaveUploadedFile saves an uploaded file (photo or GIF) and returns content type & file path
func (db *appdbimpl) SaveUploadedFile(file io.Reader, header *multipart.FileHeader, userID string) (string, string, error) {
	// Allowed file types
	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]string{
		".jpg":  "photo",
		".jpeg": "photo",
		".png":  "photo",
		".gif":  "gif",
	}

	contentType, ok := allowedExts[fileExt]
	if !ok {
		return "", "", errors.New("invalid file type")
	}

	// Ensure uploads directory exists
	uploadDir := "webui/public/uploads"
	_ = os.MkdirAll(uploadDir, os.ModePerm)

	// Generate unique filename (userID + timestamp)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	fileName := userID + "_" + timestamp + fileExt
	filePath := filepath.Join(uploadDir, fileName)

	// Save file
	out, err := os.Create(filePath)
	if err != nil {
		return "", "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", "", err
	}

	// Return content type and file path
	return contentType, "/uploads/" + fileName, nil
}

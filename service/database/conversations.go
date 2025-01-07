package database
import "database/sql"

func (db *appdbimpl) GetMyConversations_db(token string) (conversations []Conversation, err error) {
	query := `SELECT 
			conversations.id, conversations.lastconvo, conversations.is_group, conversations.photo, conversations.name
		FROM 
			usertokens
		JOIN 
			users ON usertokens.user_id = users.id
		JOIN 
			convmembers ON users.id = convmembers.user_id
		JOIN 
			conversations ON convmembers.conversation_id = conversations.id
		WHERE 
			usertokens.token = ?;`
	
	convos, err := db.c.Query(query, token)
	if err != nil {
		return
	} 
	for convos.Next() {
		var convo Conversation
		err = convos.Scan(&convo.ID, &convo.LastConvo, &convo.IsGroup, &convo.Photo, &convo.Name)
		if err != nil {
			return
		} 
		conversations = append(conversations, convo)

	}
	return conversations, err

	
}
// hw return list of conversations, user setphoto, user setusername patch requests, how to accept array from body


func (db *appdbimpl) CreateConversation_db(isGroup bool, name string, photo string) (conversation Conversation, err error) {
	query := `INSERT INTO conversations ( lastconvo, is_group, name, photo ) values (current_timestamp, ?, ?, ?) returning id`
	err = db.c.QueryRow(query, isGroup, name, photo).Scan(&conversation.ID) // retirive one row QueryRow, Exec executer 
	if err != nil && err != sql.ErrNoRows{
		return
	} else if err == sql.ErrNoRows {
		err = nil
		return
	} else {
		return conversation, nil
	}

}

func (db *appdbimpl) AddUsersToConversation(user_id int, conversation_id int) (err error) {
	query := `INSERT INTO convmembers (user_id, conversation_id) values (?, ?)`
	_, err = db.c.Exec(query, user_id, conversation_id)
	if err != nil {
		return
	} 
	return err
}

func (db *appdbimpl) GetConversationById(conversationId int) (conversation Conversation, err error) {
    query := "SELECT id, lastconvo, is_group, photo, name FROM conversations WHERE id = ?;"
    err = db.c.QueryRow(query, conversationId).Scan(
        &conversation.ID,
        &conversation.LastConvo,
        &conversation.IsGroup,
        &conversation.Photo,
        &conversation.Name,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return conversation, nil // Return an empty conversation if no rows are found
        }
        return conversation, err
    }
    return conversation, nil
}


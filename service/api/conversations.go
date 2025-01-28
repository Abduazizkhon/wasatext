package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Abduazizkhon/wasatext/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, context *reqcontext.RequestContext) {
	// Log the request
	context.Logger.Info("Request received: method=%s, path=%s", r.Method, r.URL.Path)

	// Retrieve user_id from the path parameters
	userID := ps.ByName("id")
	if userID == "" {
		context.Logger.Error("User ID is required in the path")
		http.Error(w, "User ID is required in the path", http.StatusBadRequest)
		return
	}
	context.Logger.Infof("Extracted user_id: %s", userID)

	// Fetch conversations from the database
	conversations, err := rt.db.GetMyConversations_db(userID)
	if err != nil {
		context.Logger.WithError(err).Error("Error fetching conversations")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	context.Logger.Infof("Fetched %d conversations", len(conversations))

	// Respond with the list of conversations
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(conversations)
}
func (rt *_router) sendMessageFirst(w http.ResponseWriter, r *http.Request, ps httprouter.Params, context *reqcontext.RequestContext) {
	// Log the request
	context.Logger.Info("Request received: method=%s, path=%s", r.Method, r.URL.Path)

	// Retrieve sender_id from the path parameters
	senderID := ps.ByName("id")
	if senderID == "" {
		context.Logger.Error("Sender ID is required in the path")
		http.Error(w, "Sender ID is required in the path", http.StatusBadRequest)
		return
	}
	context.Logger.Infof("Extracted sender_id: %s", senderID)

	// Parse the request body
	var input struct {
		RecipientID string `json:"recipient_id"` // Recipient's UUID
		Content     string `json:"content"`      // Message content
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.RecipientID == "" || input.Content == "" {
		context.Logger.WithError(err).Error("Invalid input: recipient_id and content are required")
		http.Error(w, "Invalid input: recipient_id and content are required", http.StatusBadRequest)
		return
	}
	context.Logger.Infof("Parsed input: recipient_id=%s, content=%s", input.RecipientID, input.Content)

	// Validate sender
	sender, err := rt.db.GetUserId(senderID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.Logger.WithError(err).Error("Sender not found")
			http.Error(w, "Sender not found", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error fetching sender")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	context.Logger.Infof("Sender found: %+v", sender)

	// Validate recipient
	recipient, err := rt.db.GetUserId(input.RecipientID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.Logger.WithError(err).Error("Recipient not found")
			http.Error(w, "Recipient not found", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error fetching recipient")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	context.Logger.Infof("Recipient found: %+v", recipient)

	// Check if a conversation already exists between the sender and recipient
	exists, err := rt.db.ConversationExists(senderID, input.RecipientID)
	if err != nil {
		context.Logger.WithError(err).Error("Error checking for existing conversation")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if exists {
		context.Logger.Error("Conversation already exists between the users")
		http.Error(w, "A conversation already exists between these users", http.StatusConflict)
		return
	}

	// Create a new conversation
	newConvo, err := rt.db.CreateConversation_db(false, recipient.Username, recipient.Photo.String)
	if err != nil {
		context.Logger.WithError(err).Error("Error creating conversation")
		http.Error(w, "Error creating conversation", http.StatusInternalServerError)
		return
	}
	context.Logger.Infof("Conversation created: %+v", newConvo)

	// Add both users to the conversation
	err = rt.db.AddUsersToConversation(sender.ID, newConvo.ID)
	if err != nil {
		context.Logger.WithError(err).Error("Error adding sender to conversation")
		http.Error(w, "Error adding sender to conversation", http.StatusInternalServerError)
		return
	}
	context.Logger.Info("Sender added to conversation")

	err = rt.db.AddUsersToConversation(recipient.ID, newConvo.ID)
	if err != nil {
		context.Logger.WithError(err).Error("Error adding recipient to conversation")
		http.Error(w, "Error adding recipient to conversation", http.StatusInternalServerError)
		return
	}
	context.Logger.Info("Recipient added to conversation")

	// Send the first message
	err = rt.db.SendMessage(newConvo.ID, sender.ID, input.Content)
	if err != nil {
		context.Logger.WithError(err).Error("Error sending message")
		http.Error(w, "Error sending message", http.StatusInternalServerError)
		return
	}
	context.Logger.Info("Message sent successfully")

	// Respond with success and conversation ID
	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"message": "Message sent successfully",
		"c_id":    newConvo.ID, // Return the new conversation ID
	}
	context.Logger.Infof("Sending response: %+v", response)
	_ = json.NewEncoder(w).Encode(response)
}

// update username of a user. Also change the name of that user in all convos

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, context reqcontext.RequestContext) {
	conversationId, err := strconv.Atoi(ps.ByName("id"))
	if err != nil || conversationId <= 0 {
		rt.baseLogger.WithError(err).Error("Invalid conversation ID")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid conversation ID"})
		return
	}

	conversation, err := rt.db.GetConversationById(conversationId)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Failed to fetch conversation")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch conversation"})
		return
	}

	if conversation.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Conversation not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(conversation)
}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, context *reqcontext.RequestContext) {
	// Log the request
	context.Logger.Info("Request received: method=%s, path=%s", r.Method, r.URL.Path)

	// Step 1: Extract conversation ID from the path
	conversationIDStr := ps.ByName("c_id")
	conversationID, err := strconv.Atoi(conversationIDStr)
	if err != nil || conversationID <= 0 {
		context.Logger.WithError(err).Error("Invalid conversation ID")
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}
	context.Logger.Infof("Extracted conversation_id: %d", conversationID)

	// Step 2: Extract sender ID from the request context (authenticated user)
	senderID := context.UserID
	if senderID == "" {
		context.Logger.Error("Sender ID is required")
		http.Error(w, "Sender ID is required", http.StatusUnauthorized)
		return
	}
	context.Logger.Infof("Extracted sender_id: %s", senderID)

	// Step 3: Parse the request body to get the message content
	var input struct {
		Content string `json:"content"` // Message content
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Content == "" {
		context.Logger.WithError(err).Error("Invalid input: content is required")
		http.Error(w, "Invalid input: content is required", http.StatusBadRequest)
		return
	}
	context.Logger.Infof("Parsed input: content=%s", input.Content)

	// Step 4: Validate if the sender is part of the conversation
	isMember, err := rt.db.IsUserInConversation(senderID, conversationID)
	if err != nil {
		context.Logger.WithError(err).Error("Error checking if user is in conversation")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		context.Logger.Error("Sender is not part of the conversation")
		http.Error(w, "Sender is not part of the conversation", http.StatusForbidden)
		return
	}

	// Step 5: Send the message to the database
	err = rt.db.SendMessageFull(conversationID, senderID, input.Content)
	if err != nil {
		context.Logger.WithError(err).Error("Error sending message")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	context.Logger.Info("Message sent successfully")

	// Step 6: Respond with success
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Message sent successfully",
	})
}

func (rt *_router) getMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, context *reqcontext.RequestContext) {
	// Extract conversation ID from the path parameters
	conversationID, err := strconv.Atoi(ps.ByName("c_id"))
	if err != nil || conversationID <= 0 {
		context.Logger.WithError(err).Error("Invalid conversation ID")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid conversation ID"})
		return
	}

	// Fetch the conversation details
	conversation, err := rt.db.GetConversationById(conversationID)
	if err != nil {
		context.Logger.WithError(err).Error("Failed to fetch conversation")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch conversation"})
		return
	}

	if conversation.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Conversation not found"})
		return
	}

	// Fetch all messages in the conversation
	messages, err := rt.db.GetMessagesByConversationId(conversationID)
	if err != nil {
		context.Logger.WithError(err).Error("Failed to fetch messages")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch messages"})
		return
	}

	// Prepare the response
	response := map[string]interface{}{
		"conversation": conversation,
		"messages":     messages,
	}

	// Respond with the conversation and messages
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, context *reqcontext.RequestContext) {
	// Log the request
	context.Logger.Info("Request received: method=%s, path=%s", r.Method, r.URL.Path)

	// Step 1: Extract conversation ID and message ID from the path parameters
	conversationIDStr := ps.ByName("c_id")
	conversationID, err := strconv.Atoi(conversationIDStr)
	if err != nil || conversationID <= 0 {
		context.Logger.WithError(err).Error("Invalid conversation ID")
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	messageIDStr := ps.ByName("m_id")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil || messageID <= 0 {
		context.Logger.WithError(err).Error("Invalid message ID")
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	// Step 2: Extract the user ID from the request context (authenticated user)
	userID := context.UserID
	if userID == "" {
		context.Logger.Error("User not authenticated")
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Step 3: Validate if the user is part of the conversation
	isMember, err := rt.db.IsUserInConversation(userID, conversationID)
	if err != nil {
		context.Logger.WithError(err).Error("Error checking if user is in conversation")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		context.Logger.Error("User is not part of the conversation")
		http.Error(w, "User is not part of the conversation", http.StatusForbidden)
		return
	}

	// Step 4: Validate if the message belongs to the user
	isMessageOwner, err := rt.db.IsMessageOwner(userID, messageID)
	if err != nil {
		context.Logger.WithError(err).Error("Error checking if user owns the message")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isMessageOwner {
		context.Logger.Error("User does not own the message")
		http.Error(w, "User does not own the message", http.StatusForbidden)
		return
	}

	// Step 5: Delete the message
	err = rt.db.DeleteMessage(messageID)
	if err != nil {
		context.Logger.WithError(err).Error("Error deleting message")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Step 6: Respond with success
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Message deleted successfully",
	})
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, context *reqcontext.RequestContext) {
	// Log the request
	context.Logger.Info("Request received: method=%s, path=%s", r.Method, r.URL.Path)

	// Step 1: Extract source conversation ID, message ID, and target conversation ID from the path parameters
	sourceConversationIDStr := ps.ByName("conversation_id")
	sourceConversationID, err := strconv.Atoi(sourceConversationIDStr)
	if err != nil || sourceConversationID <= 0 {
		context.Logger.WithError(err).Error("Invalid source conversation ID")
		http.Error(w, "Invalid source conversation ID", http.StatusBadRequest)
		return
	}

	messageIDStr := ps.ByName("message_id")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil || messageID <= 0 {
		context.Logger.WithError(err).Error("Invalid message ID")
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	targetConversationIDStr := ps.ByName("target_conversation_id")
	targetConversationID, err := strconv.Atoi(targetConversationIDStr)
	if err != nil || targetConversationID <= 0 {
		context.Logger.WithError(err).Error("Invalid target conversation ID")
		http.Error(w, "Invalid target conversation ID", http.StatusBadRequest)
		return
	}

	// Step 2: Extract the user ID from the request context (authenticated user)
	userID := context.UserID
	if userID == "" {
		context.Logger.Error("User not authenticated")
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Step 3: Validate if the user is part of both the source and target conversations
	isMemberOfSource, err := rt.db.IsUserInConversation(userID, sourceConversationID)
	if err != nil {
		context.Logger.WithError(err).Error("Error checking if user is in source conversation")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isMemberOfSource {
		context.Logger.Error("User is not part of the source conversation")
		http.Error(w, "User is not part of the source conversation", http.StatusForbidden)
		return
	}

	isMemberOfTarget, err := rt.db.IsUserInConversation(userID, targetConversationID)
	if err != nil {
		context.Logger.WithError(err).Error("Error checking if user is in target conversation")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isMemberOfTarget {
		context.Logger.Error("User is not part of the target conversation")
		http.Error(w, "User is not part of the target conversation", http.StatusForbidden)
		return
	}

	// Step 4: Fetch the message content from the source conversation
	messageContent, err := rt.db.GetMessageContent(messageID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.Logger.WithError(err).Error("Message not found")
			http.Error(w, "Message not found", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error fetching message content")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Step 5: Forward the message to the target conversation with the user as the sender
	err = rt.db.ForwardMessage(targetConversationID, userID, messageContent)
	if err != nil {
		context.Logger.WithError(err).Error("Error forwarding message")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Step 6: Respond with success
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Message forwarded successfully",
	})
}

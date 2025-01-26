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


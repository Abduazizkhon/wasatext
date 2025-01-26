package api

import (
	"encoding/json"

	"net/http"

	// "github.com/Abduazizkhon/wasatext/service/api/reqcontext"
	// "github.com/gofrs/uuid"
	"github.com/Abduazizkhon/wasatext/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"

	// "github.com/sirupsen/logrus"
	"database/sql"

)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, context *reqcontext.RequestContext) {
	// Parse input
	var input struct {
		Username string `json:"username"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Username == "" {
		http.Error(w, "Invalid input: username is required", http.StatusBadRequest)
		return
	}

	// Check if the user exists
	user, err := rt.db.GetUser(input.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			// User doesn't exist, create them
			user, err = rt.db.CreateUser(input.Username)
			if err != nil {
				rt.baseLogger.WithError(err).Error("Failed to create user")
				http.Error(w, "Internal server error: failed to create user", http.StatusInternalServerError)
				return
			}
		} else {
			// Log unexpected errors
			rt.baseLogger.WithError(err).Error("Unexpected error fetching user")
			http.Error(w, "Internal server error: unexpected error", http.StatusInternalServerError)
			return
		}
	}

	// Respond with the user data and user ID (token)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user,
		"token": user.ID, // The user ID is the token
	})
}
func (rt *_router) logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params, context reqcontext.RequestContext) {
	// Parse the request body to extract the user ID (UUID)
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.ID == "" {
		http.Error(w, "Invalid request body or missing user ID", http.StatusBadRequest)
		return
	}

	// Validate if the user exists in the database
	_, err = rt.db.GetUserId(user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to validate user", http.StatusInternalServerError)
		return
	}

	// Perform any additional logout logic here (if needed)

	// Send a success response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Logout successful"))
}

// -----------------

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, context *reqcontext.RequestContext) {
	// Extract the user ID from the context
	userID := context.UserID
	if userID == "" {
		rt.baseLogger.Error("User not authenticated")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "User not authenticated"})
		return
	}

	// Log the user ID for debugging
	rt.baseLogger.Infof("setMyUserName: userID=%s", userID)

	// Parse input
	var input struct {
		NewName string `json:"newname"` // New username
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Invalid input")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
		return
	}
	rt.baseLogger.Infof("setMyUserName: newname=%s", input.NewName)

	if input.NewName == "" {
		rt.baseLogger.Error("setMyUserName: New username is required")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "New username is required"})
		return
	}

	// Check if the new username already exists
	existName, err := rt.db.GetUser(input.NewName)
	if err != nil && err != sql.ErrNoRows {
		rt.baseLogger.WithError(err).Error("setMyUserName: Error checking existing username")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	if existName.Username != "" {
		rt.baseLogger.Error("setMyUserName: Username already exists")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Username already exists"})
		return
	}

	// Update username in the database
	err = rt.db.UpdateUserName(userID, input.NewName)
	if err != nil {
		rt.baseLogger.WithError(err).Error("setMyUserName: Failed to update username")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update username"})
		return
	}

	// Success response
	rt.baseLogger.Info("setMyUserName: Username updated successfully")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Username updated successfully"})
}
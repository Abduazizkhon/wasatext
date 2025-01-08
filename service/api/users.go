package api

import (
	"encoding/json"


	"net/http"

	"github.com/Abduazizkhon/wasatext/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)


func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user) // r is request
	if err != nil {
		w.WriteHeader(400)
		return
	}
	user_obj, err := rt.db.GetUser(user.Username)

	if err != nil {
		http.Error(w, "Error decoding request body: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})
	
		if user_obj.ID == 0 {
			user_obj, err := rt.db.CreateUser(user.Username)
			if err != nil {
				rt.baseLogger.WithError(err).Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)				
				return
			} 
			user_obj.Token = reqUUID.String()
			err = rt.db.SetToken(user_obj.ID, user_obj.Token)
			if err != nil {
				rt.baseLogger.WithError(err).Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)				
				return
			} 
			w.WriteHeader(201)
			w.Header().Set("content-type", "application/json")
			_= json.NewEncoder(w).Encode(user_obj)
			
		} else{
			user_obj.Token = reqUUID.String()
			err = rt.db.SetToken(user_obj.ID, user_obj.Token)
			if err != nil {
				rt.baseLogger.WithError(err).Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)				
				return
			} 
			w.WriteHeader(200)
			w.Header().Set("content-type", "application/json")
			_= json.NewEncoder(w).Encode(user_obj)

		}
		
		
	}
}

func (rt *_router) logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var token UserToken
	err := json.NewDecoder(r.Body).Decode(&token) // r is request
	if err != nil {
		w.WriteHeader(400)
		return
	}
	err = rt.db.DeleteToken(token.Token)
	if err != nil {
		// rt.baseLogger.WithError(err).Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)				
		return
	} 
	w.WriteHeader(200)
	w.Header().Set("content-type", "text/plaint")
	_= json.NewEncoder(w).Encode("Logout is done")
	
	
}
// -----------------

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    // Parse input
    var input struct {
        Token       string    `json:"token"`
        NewName  string `json:"newname"`
    }

    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil {
        rt.baseLogger.WithError(err).Error("Invalid input")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
        return
    }

    if input.Token == "" || input.NewName == "" {
        rt.baseLogger.Error("Missing required fields")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(map[string]string{"error": "Missing required fields"})
        return
    }

	existName, err := rt.db.GetUser(input.NewName)
	if existName.Username != "" {
		w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(map[string]string{"error": "Username already exists"})
        return
	}
	if err != nil {
		// rt.baseLogger.WithError(err).Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)				
		return
	} 




    // Update username in the database
	userid, err := rt.db.GetUserId(input.Token)
	if err != nil {
		// rt.baseLogger.WithError(err).Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)				
		return
	} 
    err = rt.db.UpdateUserName(userid.User_id, input.NewName)
    if err != nil {
        rt.baseLogger.WithError(err).Error("Failed to update username")
        w.WriteHeader(http.StatusInternalServerError)
        _ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update username"})
        return
    }

    // Success response
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(map[string]string{"message": "Username updated successfully"})
}


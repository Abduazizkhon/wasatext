package api
import (
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) GetMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var token UserToken
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	convos, err := rt.db.GetMyConversations_db(token.Token)
	w.WriteHeader(200)
	w.Header().Set("content-type", "application/json")
	_= json.NewEncoder(w).Encode(convos)
	if err != nil {
		rt.baseLogger.WithError(err).Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)				
		return
	} 
	
}

func (rt *_router) SendFirstMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	seconduser, err := rt.db.GetUser(user.Username)
	if err != nil {
		rt.baseLogger.WithError(err).Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)				
		return
	} 
	if seconduser.ID == 0 {
		w.WriteHeader(400)
		w.Header().Set("content-type", "text/plaint")
		_= json.NewEncoder(w).Encode("User does not exist")				
		return
	}
	usertoken, err := rt.db.GetUserId(user.Token)
	if usertoken.User_id == 0 {
		w.WriteHeader(403)
		w.Header().Set("content-type", "text/plaint")
		_= json.NewEncoder(w).Encode("User is not authorised. Do the login")				
		return
	}
	if err != nil {
		rt.baseLogger.WithError(err).Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)				
		return
	} 
	convos, err := rt.db.GetMyConversations_db(user.Token)
	if err != nil {
		rt.baseLogger.WithError(err).Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)				
		return
	} 
	for i := 0; i < len(convos); i++ {
		if convos[i].Name == seconduser.Username {
			w.WriteHeader(400)
			w.Header().Set("content-type", "text/plaint")
			_= json.NewEncoder(w).Encode("Such chat is already exists")				
			return
		}
    } 
	newconvo, err := rt.db.CreateConversation_db(false, seconduser.Username, seconduser.Photo.String)
	if err != nil {
		rt.baseLogger.WithError(err).Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)				
		return
	} 
	err = rt.db.AddUsersToConversation(usertoken.User_id, newconvo.ID)
	if err != nil {
		rt.baseLogger.WithError(err).Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)				
		return
	} 

	err = rt.db.AddUsersToConversation(seconduser.ID, newconvo.ID)
	if err != nil {
		rt.baseLogger.WithError(err).Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)				
		return
	} 

	w.WriteHeader(200)
	w.Header().Set("content-type", "text/plaint")
	_= json.NewEncoder(w).Encode("Chat is created")				
}

// update username of a user. Also change the name of that user in all convos

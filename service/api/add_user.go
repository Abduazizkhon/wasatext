package api

import (
 "encoding/json"
 "github.com/gofrs/uuid"
 "github.com/julienschmidt/httprouter"
 "net/http"
)

func (rt *_router) addUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
 if r.Method != http.MethodPut {
  http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
  return
 }
 w.Header().Set("content-type", "application/json")

 var username UserName
 err := json.NewDecoder(r.Body).Decode(&username)
 if err != nil {
  w.WriteHeader(http.StatusBadRequest)
  return
 }
 println(username.UserName)
 println("Username Received")
 if username.UserName == "" {
  w.WriteHeader(http.StatusBadRequest)
  return
 }

 var currentUser User
 currentUser = User{Id: uuid.Must(uuid.NewV4()), UserName: username.UserName}
 println("Username " + currentUser.UserName)
 println(currentUser.Id.String())
 err = rt.db.AddUsername(currentUser.Id, currentUser.UserName)
 if err != nil {
  print("Error adding user")
 }
 //UsersDict[currentUser.Id] = len(Users)
 //Users = append(Users, username)
 w.WriteHeader(http.StatusCreated)
 type UserResponse struct {
  Id       uuid.UUID
  Username string
 }
 json.NewEncoder(w).Encode(UserResponse{
  Id:       currentUser.Id,
  Username: currentUser.UserName,
 })
}
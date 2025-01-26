package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	// rt.router.GET("/liveness/", rt.wrap(rt.liveness))
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	// rt.router.POST("/logout", rt.wrap(rt.logout))
	rt.router.PUT("/users/me/username", rt.wrap(rt.setMyUserName))
	rt.router.GET("/users/:id/conversations", rt.wrap(rt.getMyConversations))
	rt.router.POST("/users/:id/conversations/first-message", rt.wrap(rt.sendMessageFirst))	
	// rt.router.POST("/conversations/:c_id/messages", rt.wrap(rt.sendMessage))// Send message to an existing conversation
	// rt.router.GET("/users/:id/conversations/:c_id", rt.getConversation)
	// :conversation

	return rt.router
}
// use context in every(not dologin) api 
// after I iterate over rows or use them  in database I need to
// if err = rows.Err(); err != nil {
// } 
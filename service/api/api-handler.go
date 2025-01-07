package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/liveness/", rt.liveness)
	rt.router.POST("/session", rt.doLogin)
	rt.router.POST("/logout", rt.logout)
	rt.router.GET("/users/:id/conversations", rt.getMyConversations)
	rt.router.POST("/users/:id/conversations/:c_id/messages", rt.sendMessage)
	rt.router.PUT("/users/me/username", rt.setMyUserName)
	rt.router.GET("/users/:id/conversations/:c_id", rt.getConversation)
// :conversation


	return rt.router
}

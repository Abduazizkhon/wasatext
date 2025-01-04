package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/liveness/", rt.liveness)
	rt.router.POST("/login/", rt.login)
	rt.router.POST("/logout", rt.logout)
	rt.router.GET("/users/myconversations/", rt.GetMyConversations)
	rt.router.POST("/users/myconversations/newconvo", rt.SendFirstMessage)
// :conversation


	return rt.router
}

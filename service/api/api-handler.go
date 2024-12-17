package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	// my code
	rt.router.PUT("/users/me/username", rt.setUsername) //setUsername is the name of the function that I need to define
	rt.router.POST("/session", rt.doLogin)
	// Conversations
	rt.router.GET("/users/{id}/conversations", rt.getConvo)
	rt.router.GET("/users/{id}/conversations/{c_id}", rt.getConversation)
	// Messages
	rt.router.GET("/users/{id}/conversations/{c_id}/messages", rt.getMessagesInConversation)
	rt.router.POST("/users/{id}/conversations/{c_id}/messages", rt.sendMessage)
	rt.router.POST("/users/{id}/conversations/{c_id}/messages/{m_id}/forward", rt.forwardMessage)
	rt.router.POST("/users/{id}/conversations/{c_id}/messages/{m_id}/comments", rt.commentMessage)
	rt.router.DELETE("/users/{id}/conversations/{c_id}/messages/{m_id}/comments/{commentId}", rt.uncommentMessage)
	rt.router.DELETE("/users/{id}/conversations/{c_id}/messages/{m_id}", rt.deleteMessage)
	// Groups Routes
	rt.router.POST("/users/{id}/conversations/{c_id}/groups/{g_id}/members", rt.addToGroup)
	rt.router.POST("/users/{id}/conversations/{c_id}/groups/{g_id}/leave", rt.leaveGroup)
	rt.router.PUT("/users/{id}/conversations/{c_id}/groups/{g_id}/name", rt.setGroupName)
	rt.router.PUT("/users/{id}/conversations/{c_id}/groups/{g_id}/photo", rt.setGroupPhoto)



	// ----------

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}

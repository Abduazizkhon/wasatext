package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	// rt.router.GET("/liveness/", rt.wrap(rt.liveness))
	rt.router.POST("/session", rt.wrap(rt.doLogin)) // done
	rt.router.ServeFiles("/uploads/*filepath", http.Dir("webui/public/uploads"))

	// rt.router.POST("/logout", rt.wrap(rt.logout))
	rt.router.PUT("/users/me/username", rt.wrap(rt.setMyUserName)) // done
	rt.router.PUT("/users/me/photo", rt.wrap(rt.setMyPhoto))       // done
	rt.router.GET("/users/:id/conversations", rt.wrap(rt.getMyConversations)) // done
	rt.router.POST("/users/:id/conversations/first-message", rt.wrap(rt.sendMessageFirst)) // done
	rt.router.POST("/conversations/:conversation_id/messages", rt.wrap(rt.sendMessage)) // done
	rt.router.GET("/conversations/:c_id", rt.wrap(rt.getMessages)) // done
	rt.router.DELETE("/conversations/:conversation_id/messages/:message_id", rt.wrap(rt.deleteMessage)) // done
	rt.router.POST("/conversations/:conversation_id/messages/:message_id/forward/:target_conversation_id", rt.wrap(rt.forwardMessage))
	rt.router.POST("/groups", rt.wrap(rt.createGroup)) // done
	rt.router.POST("/groups/:c_id/members", rt.wrap(rt.addToGroup)) // done
	rt.router.DELETE("/groups/:c_id/leave", rt.wrap(rt.leaveGroup)) // done
	rt.router.PUT("/groups/:c_id/name", rt.wrap(rt.setGroupName))
	rt.router.PUT("/conversations/:c_id/set-group-photo", rt.wrap(rt.setGroupPhoto))
	rt.router.POST("/conversations/:conversation_id/messages/:message_id/comments", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/conversations/:conversation_id/messages/:message_id/comments/:comment_id", rt.wrap(rt.uncommentMessage))
	rt.router.GET("/users/:id", rt.wrap(rt.getUser)) // ✅ Add this route
	

	// rt.router.POST("/conversations/:c_id/messages", rt.wrap(rt.sendMessage))// Send message to an existing conversation
	// rt.router.GET("/users/:id/conversations/:c_id", rt.getConversation)
	// :conversation

	return rt.router
}

// use context in every(not dologin) api
// after I iterate over rows or use them  in database I need to
// if err = rows.Err(); err != nil {
// }

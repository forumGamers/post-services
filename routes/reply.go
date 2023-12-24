package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/post-services/controller"
	md "github.com/post-services/middlewares"
)

func (r routes) replyRoutes(rg *gin.RouterGroup, rc controller.ReplyController, middleware md.Middleware) {
	uri := rg.Group("/reply")

	uri.Use(middleware.SetContexts)
	uri.Use(middleware.Authentication)
	uri.POST("/:commentId", rc.AddReply)
	uri.DELETE("/:commentId/:replyId", rc.DeleteReply)
}

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/post-services/controller"
	md "github.com/post-services/middlewares"
)

func (r routes) replyRoutes(rg *gin.RouterGroup, rc controller.ReplyController) {
	uri := rg.Group("/reply")

	uri.POST("/:commentId", md.Authentication, rc.AddReply)
}

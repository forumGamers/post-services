package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/post-services/controller"
	md "github.com/post-services/middlewares"
)

func (r routes) commentRoutes(rg *gin.RouterGroup, pc controller.CommentController) {
	uri := rg.Group("/comment")

	uri.POST("/:postId", md.Authentication, pc.CreateComment)
}

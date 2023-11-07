package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/post-services/controller"
	md "github.com/post-services/middlewares"
)

func (r routes) commentRoutes(rg *gin.RouterGroup, cc controller.CommentController) {
	uri := rg.Group("/comment")

	uri.Use(md.SetContext)
	uri.Use(md.Authentication)
	uri.POST("/bulk", cc.BulkComment)
	uri.POST("/:postId", cc.CreateComment)
	uri.DELETE("/:commentId", cc.DeleteComment)
}

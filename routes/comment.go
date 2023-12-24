package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/post-services/controller"
	md "github.com/post-services/middlewares"
)

func (r routes) commentRoutes(rg *gin.RouterGroup, cc controller.CommentController, middleware md.Middleware) {
	uri := rg.Group("/comment")

	uri.Use(middleware.SetContexts)
	uri.Use(middleware.Authentication)
	uri.POST("/bulk", cc.BulkComment)
	uri.POST("/:postId", cc.CreateComment)
	uri.DELETE("/:commentId", cc.DeleteComment)
}

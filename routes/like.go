package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/post-services/controller"
	md "github.com/post-services/middlewares"
)

func (r routes) likeRoutes(rg *gin.RouterGroup, lc controller.LikeController, middleware md.Middleware) {
	uri := rg.Group("/like")

	uri.Use(middleware.SetContexts)
	uri.Use(middleware.Authentication)
	uri.POST("/bulk", lc.BulkLikes)
	uri.POST("/:postId", lc.LikePost)
	uri.DELETE("/:postId", lc.UnlikePost)
}

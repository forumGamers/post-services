package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/post-services/controller"
	md "github.com/post-services/middlewares"
)

func (r routes) postRoutes(rg *gin.RouterGroup, pc controller.PostController, middleware md.Middleware) {
	uri := rg.Group("/post")

	uri.Use(middleware.SetContexts)
	uri.Use(middleware.Authentication)
	uri.POST("/", pc.CreatePost)
	uri.POST("/bulk", pc.BulkCreatePost)
	uri.DELETE("/:postId", pc.DeletePost)
}

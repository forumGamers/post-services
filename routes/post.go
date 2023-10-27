package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/post-services/controller"
	md "github.com/post-services/middlewares"
)

func (r routes) postRoutes(rg *gin.RouterGroup, pc controller.PostController) {
	uri := rg.Group("/post")

	uri.Use(md.SetContext)
	uri.POST("/", md.Authentication, pc.CreatePost)
	uri.POST("/bulk", md.Authentication, pc.BulkCreatePost)
	uri.DELETE("/:postId", md.Authentication, pc.DeletePost)
}

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/post-services/controller"
	md "github.com/post-services/middlewares"
)

func (r routes) likeRoutes(rg *gin.RouterGroup,lc controller.LikeController) {
	uri := rg.Group("/like")

	uri.POST("/:postId",md.Authentication,lc.LikePost)
}
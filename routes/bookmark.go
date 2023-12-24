package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/post-services/controller"
	md "github.com/post-services/middlewares"
)

func (r routes) bookmarkRoute(rg *gin.RouterGroup, bc controller.BookmarkController, middleware md.Middleware) {
	uri := rg.Group("/bookmark")

	uri.Use(middleware.SetContexts)
	uri.Use(middleware.Authentication)
	uri.POST("/:postId", bc.CreateBookmark)
}

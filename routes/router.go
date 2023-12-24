package routes

import (
	"os"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	c "github.com/post-services/controller"
	h "github.com/post-services/helper"
	md "github.com/post-services/middlewares"
)

type routes struct {
	router *gin.Engine
}

// add limiter
func NewRouter(
	middleware md.Middleware,
	post c.PostController,
	like c.LikeController,
	comment c.CommentController,
	reply c.ReplyController,
) {
	h.PanicIfError(godotenv.Load())

	r := routes{router: gin.Default()}

	r.router.Use(middleware.CheckOrigin)
	r.router.Use(middleware.Cors())
	r.router.Use(logger.SetLogger())

	groupRoutes := r.router.Group("/api/v1")

	r.postRoutes(groupRoutes, post, middleware)
	r.likeRoutes(groupRoutes, like, middleware)
	r.commentRoutes(groupRoutes, comment, middleware)
	r.replyRoutes(groupRoutes, reply, middleware)

	port := os.Getenv("PORT")

	if port == "" {
		port = "4300"
	}

	r.router.Run(":" + port)
}

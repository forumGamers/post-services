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
	post c.PostController,
	like c.LikeController,
	comment c.CommentController,
	reply c.ReplyController,
) {
	h.PanicIfError(godotenv.Load())

	r := routes{router: gin.Default()}

	groupRoutes := r.router.Group("/api")

	r.router.Use(md.CheckOrigin)
	r.router.Use(md.Cors())
	r.router.Use(logger.SetLogger())
	r.postRoutes(groupRoutes, post)
	r.likeRoutes(groupRoutes, like)
	r.commentRoutes(groupRoutes, comment)
	r.replyRoutes(groupRoutes, reply)

	port := os.Getenv("PORT")

	if port == "" {
		port = "4300"
	}

	r.router.Run(":" + port)
}

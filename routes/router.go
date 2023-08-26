package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	c "github.com/post-services/controller"
	h "github.com/post-services/helper"
)

type routes struct {
	router *gin.Engine
}

func NewRouter(
	post	c.PostController,
) {
	h.PanicIfError(godotenv.Load())

	r := routes { router: gin.Default() }

	groupRoutes := r.router.Group("/api")

	r.postRoutes(groupRoutes,post)

	port := os.Getenv("PORT")

	if port == "" {
		port = "4300"
	}

	r.router.Run(":"+port)
}
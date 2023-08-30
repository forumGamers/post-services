package main

import (
	"github.com/joho/godotenv"
	cfg "github.com/post-services/config"
	c "github.com/post-services/controller"
	h "github.com/post-services/helper"
	l "github.com/post-services/pkg/like"
	p "github.com/post-services/pkg/post"
	"github.com/post-services/routes"
	tp "github.com/post-services/third-party"
	v "github.com/post-services/validations"
)

func main() {
	h.PanicIfError(godotenv.Load())
	cfg.Connection()
	validate := v.GetValidator()
	imageKit := tp.ImageKitConnection()

	postRepo := p.NewPostRepo()
	postController := c.NewPostController(p.NewPostService(postRepo,validate,imageKit),postRepo)

	likeRepo := l.NewLikeRepo()
	likeController := c.NewLikeController(l.NewLikeService(likeRepo,validate),likeRepo)

	routes.NewRouter(
		postController,
		likeController,
	)
}
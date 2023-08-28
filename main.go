package main

import (
	"github.com/joho/godotenv"
	cfg "github.com/post-services/config"
	c "github.com/post-services/controller"
	h "github.com/post-services/helper"
	p "github.com/post-services/pkg/post"
	"github.com/post-services/routes"
	tp "github.com/post-services/third-party"
	v "github.com/post-services/validations"
)

func main() {
	h.PanicIfError(godotenv.Load())
	db := cfg.Connection()
	validate := v.GetValidator()
	imageKit := tp.ImageKitConnection()

	postRepo := p.NewPostRepo(db.Collection("post"))
	postService := p.NewPostService(postRepo,validate,imageKit)
	postController := c.NewPostController(postService,postRepo)

	routes.NewRouter(postController)
}
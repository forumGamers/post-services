package main

import (
	"github.com/joho/godotenv"
	cfg "github.com/post-services/config"
	c "github.com/post-services/controller"
	h "github.com/post-services/helper"
	"github.com/post-services/routes"
	// tp "github.com/post-services/third-party"
	v "github.com/post-services/validations"

	r "github.com/post-services/repository"
	s "github.com/post-services/services"
)

func main() {
	h.PanicIfError(godotenv.Load())
	db := cfg.Connection()
	validate := v.GetValidator()
	// imageKit := tp.ImageKitConnection()

	postRepo := r.NewPostRepo(db.Collection("post"))
	postService := s.NewPostService(postRepo,validate)
	postController := c.NewPostController(postService,postRepo)

	routes.NewRouter(postController)
}
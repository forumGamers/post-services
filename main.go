package main

import (
	"github.com/joho/godotenv"
	"github.com/post-services/broker"
	cfg "github.com/post-services/config"
	c "github.com/post-services/controller"
	h "github.com/post-services/helper"
	com "github.com/post-services/pkg/comment"
	l "github.com/post-services/pkg/like"
	p "github.com/post-services/pkg/post"
	r "github.com/post-services/pkg/reply"
	"github.com/post-services/pkg/share"
	"github.com/post-services/routes"
	tp "github.com/post-services/third-party"
	v "github.com/post-services/validations"
)

func main() {
	h.PanicIfError(godotenv.Load())
	cfg.Connection()
	broker.BrokerConnection()

	validate := v.GetValidator()
	imageKit := tp.ImageKitConnection()

	shareRepo := share.NewShareRepo()
	likeRepo := l.NewLikeRepo()
	commentRepo := com.NewCommentRepo()
	postRepo := p.NewPostRepo()

	postController := c.NewPostController(p.NewPostService(postRepo, validate, imageKit), postRepo, commentRepo, likeRepo, shareRepo)
	likeController := c.NewLikeController(l.NewLikeService(likeRepo, validate), likeRepo)
	commentController := c.NewCommentController(com.NewCommentService(commentRepo, validate), commentRepo)
	replyController := c.NewReplyController(r.NewReplyService(commentRepo, validate), commentRepo)

	routes.NewRouter(
		postController,
		likeController,
		commentController,
		replyController,
	)
}

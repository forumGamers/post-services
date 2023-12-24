package main

import (
	"github.com/joho/godotenv"
	"github.com/post-services/broker"
	cfg "github.com/post-services/config"
	c "github.com/post-services/controller"
	h "github.com/post-services/helper"
	md "github.com/post-services/middlewares"
	com "github.com/post-services/pkg/comment"
	l "github.com/post-services/pkg/like"
	p "github.com/post-services/pkg/post"
	"github.com/post-services/pkg/reply"
	"github.com/post-services/pkg/share"
	"github.com/post-services/routes"
	tp "github.com/post-services/third-party"
	v "github.com/post-services/validations"
	"github.com/post-services/web"
)

func main() {
	h.PanicIfError(godotenv.Load())
	cfg.Connection()
	broker.BrokerConnection()

	validate := v.GetValidator()
	imageKit := tp.ImageKitConnection()
	w := web.NewResponseWriter()
	r := web.NewRequestReader()
	middleware := md.NewMiddlewares(w, r)

	shareRepo := share.NewShareRepo()
	likeRepo := l.NewLikeRepo()
	commentRepo := com.NewCommentRepo()
	postRepo := p.NewPostRepo()

	postController := c.NewPostController(p.NewPostService(postRepo, validate, imageKit), postRepo, commentRepo, likeRepo, shareRepo, r, w)
	likeController := c.NewLikeController(l.NewLikeService(likeRepo, validate), likeRepo, r, w)
	commentController := c.NewCommentController(com.NewCommentService(commentRepo, validate), commentRepo, r, w)
	replyController := c.NewReplyController(reply.NewReplyService(commentRepo, validate), commentRepo, r, w)

	routes.NewRouter(
		middleware,
		postController,
		likeController,
		commentController,
		replyController,
	)
}

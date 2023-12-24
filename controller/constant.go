package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/post-services/pkg/bookmark"
	"github.com/post-services/pkg/comment"
	"github.com/post-services/pkg/like"
	l "github.com/post-services/pkg/like"
	"github.com/post-services/pkg/post"
	p "github.com/post-services/pkg/post"
	r "github.com/post-services/pkg/reply"
	"github.com/post-services/pkg/share"
	"github.com/post-services/web"
)

type PostController interface {
	CreatePost(c *gin.Context)
	DeletePost(c *gin.Context)
	BulkCreatePost(c *gin.Context)
}

type PostControllerImpl struct {
	web.ResponseWriter
	web.RequestReader
	Service     p.PostService
	Repo        p.PostRepo
	CommentRepo comment.CommentRepo
	LikeRepo    like.LikeRepo
	ShareRepo   share.ShareRepo
}

type ReplyController interface {
	AddReply(c *gin.Context)
	DeleteReply(c *gin.Context)
}

type ReplyControllerImpl struct {
	web.ResponseWriter
	web.RequestReader
	Service r.ReplyService
	Comment comment.CommentRepo
}

type LikeController interface {
	LikePost(c *gin.Context)
	UnlikePost(c *gin.Context)
	BulkLikes(c *gin.Context)
}

type LikeControllerImpl struct {
	web.ResponseWriter
	web.RequestReader
	Service l.LikeService
	Repo    l.LikeRepo
}

type CommentController interface {
	CreateComment(c *gin.Context)
	DeleteComment(c *gin.Context)
	BulkComment(c *gin.Context)
}

type CommentControllerImpl struct {
	web.ResponseWriter
	web.RequestReader
	Repo    comment.CommentRepo
	Service comment.CommentService
}

type BookmarkController interface {
	CreateBookmark(c *gin.Context)
}

type BookmarkControllerImpl struct {
	web.ResponseWriter
	web.RequestReader
	Repo     bookmark.BookmarkRepo
	Service  bookmark.BookmarkService
	PostRepo post.PostRepo
}

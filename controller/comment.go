package controller

import (
	"context"

	"github.com/gin-gonic/gin"
	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	b "github.com/post-services/pkg/base"
	"github.com/post-services/pkg/comment"
	"github.com/post-services/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentController interface {
	CreateComment(c *gin.Context)
	DeleteComment(c *gin.Context)
}

type CommentControllerImpl struct {
	Repo    comment.CommentRepo
	Service comment.CommentService
}

func NewCommentController(service comment.CommentService, repo comment.CommentRepo) CommentController {
	return &CommentControllerImpl{
		Repo:    repo,
		Service: service,
	}
}

func (pc *CommentControllerImpl) CreateComment(c *gin.Context) {
	postId, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		web.AbortHttp(c, h.ErrInvalidObjectId)
		return
	}

	var data web.CommentForm
	c.ShouldBind(&data)

	if err := pc.Service.ValidateComment(&data); err != nil {
		web.HttpValidationErr(c, err)
		return
	}

	var post m.Post
	if err := b.NewBaseRepo(b.GetCollection(b.Post)).FindOneById(context.Background(), postId, &post); err != nil {
		web.AbortHttp(c, err)
		return
	}

	comment := pc.Service.CreatePayload(data, postId, h.GetUser(c).Id)
	if err := pc.Repo.CreateComment(context.Background(), &comment); err != nil {
		web.AbortHttp(c, err)
		return
	}

	web.WriteResponse(c, web.WebResponse{
		Code:    201,
		Message: "success",
		Data:    comment,
	})
}

func (pc *CommentControllerImpl) DeleteComment(c *gin.Context) {
	commentId, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		web.AbortHttp(c, h.ErrInvalidObjectId)
		return
	}

	var comment m.Comment
	if err := pc.Repo.FindById(context.Background(), commentId, &comment); err != nil {
		web.AbortHttp(c, err)
		return
	}

	if err := pc.Service.AuthorizeDeleteComment(comment, h.GetUser(c)); err != nil {
		web.AbortHttp(c, err)
		return
	}

	if err := pc.Repo.DeleteOne(context.Background(), commentId); err != nil {
		web.AbortHttp(c, err)
		return
	}

	web.WriteResponse(c, web.WebResponse{
		Code:    200,
		Message: "success",
	})
}

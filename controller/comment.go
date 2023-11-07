package controller

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	br "github.com/post-services/broker"
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
	BulkComment(c *gin.Context)
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

	comment := pc.Service.CreatePayload(data, postId, h.GetUser(c).UUID)
	if err := pc.Repo.CreateComment(context.Background(), &comment); err != nil {
		web.AbortHttp(c, err)
		return
	}

	if err := br.Broker.PublishMessage(context.Background(), br.COMMENTEXCHANGE, br.NEWCOMMENTQUEUE, "application/json", br.CommentDocumment{
		Id:        comment.Id.Hex(),
		UserId:    comment.UserId,
		Text:      comment.Text,
		PostId:    comment.PostId.Hex(),
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}); err != nil {
		web.AbortHttp(c, h.BadGateway)
		return
	}

	comment.Text = h.Decryption(comment.Text)

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

	if err := br.Broker.PublishMessage(context.Background(), br.COMMENTEXCHANGE, br.DELETECOMMENTQUEUE, "application/json", br.CommentDocumment{
		Id:        comment.Id.Hex(),
		UserId:    comment.UserId,
		Text:      comment.Text,
		PostId:    comment.PostId.Hex(),
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}); err != nil {
		web.AbortHttp(c, h.BadGateway)
		return
	}

	web.WriteResponse(c, web.WebResponse{
		Code:    200,
		Message: "success",
	})
}

func (pc *CommentControllerImpl) BulkComment(c *gin.Context) {
	if h.GetStage(c) != "Development" {
		web.CustomMsgAbortHttp(c, "No Content", 204)
		return
	}

	var datas web.CommentDatas
	c.ShouldBind(&datas)

	var comments []web.CommentData
	for _, data := range datas.Datas {
		postId, _ := primitive.ObjectIDFromHex(data.PostId.Hex())
		if data.UserId != "" && postId != primitive.NilObjectID {
			data.Text = h.Encryption(data.Text)
			data.PostId = postId
			comments = append(comments, data)
		}
	}

	pc.Service.InsertManyAndBindIds(context.Background(), comments)

	var commentDocuments []br.CommentDocumment
	for _, comment := range comments {
		commentDocuments = append(commentDocuments, br.CommentDocumment{
			Id:        comment.Id.Hex(),
			UserId:    comment.UserId,
			PostId:    comment.PostId.Hex(),
			Text:      comment.Text,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	if err := br.Broker.PublishMessage(
		context.Background(),
		br.COMMENTEXCHANGE,
		br.BULKCOMMENTQUEUE,
		"application/json",
		&commentDocuments,
	); err != nil {
		web.AbortHttp(c, h.BadGateway)
		return
	}

	web.WriteResponse(c, web.WebResponse{
		Code:    201,
		Message: "success",
		Data:    comments,
	})
}

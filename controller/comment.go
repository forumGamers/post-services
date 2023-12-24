package controller

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	br "github.com/post-services/broker"
	h "github.com/post-services/helper"
	"github.com/post-services/models"
	m "github.com/post-services/models"
	b "github.com/post-services/pkg/base"
	"github.com/post-services/pkg/comment"
	"github.com/post-services/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewCommentController(service comment.CommentService, repo comment.CommentRepo, r web.RequestReader, w web.ResponseWriter) CommentController {
	return &CommentControllerImpl{w, r, repo, service}
}

func (pc *CommentControllerImpl) CreateComment(c *gin.Context) {
	postId, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		pc.AbortHttp(c, h.ErrInvalidObjectId)
		return
	}

	var data web.CommentForm
	pc.GetParams(c, &data)

	if err := pc.Service.ValidateComment(&data); err != nil {
		pc.HttpValidationErr(c, err)
		return
	}

	var post m.Post
	if err := b.NewBaseRepo(b.GetCollection(b.Post)).FindOneById(context.Background(), postId, &post); err != nil {
		pc.AbortHttp(c, err)
		return
	}

	comment := pc.Service.CreatePayload(data, postId, h.GetUser(c).UUID)
	if err := pc.Repo.CreateComment(context.Background(), &comment); err != nil {
		pc.AbortHttp(c, err)
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
		pc.AbortHttp(c, h.BadGateway)
		return
	}

	comment.Text = h.Decryption(comment.Text)

	pc.WriteResponse(c, web.WebResponse{
		Code:    201,
		Message: "success",
		Data:    comment,
	})
}

func (pc *CommentControllerImpl) DeleteComment(c *gin.Context) {
	commentId, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		pc.AbortHttp(c, h.ErrInvalidObjectId)
		return
	}

	var comment m.Comment
	if err := pc.Repo.FindById(context.Background(), commentId, &comment); err != nil {
		pc.AbortHttp(c, err)
		return
	}

	if err := pc.Service.AuthorizeDeleteComment(comment, h.GetUser(c)); err != nil {
		pc.AbortHttp(c, err)
		return
	}

	if err := pc.Repo.DeleteOne(context.Background(), commentId); err != nil {
		pc.AbortHttp(c, err)
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
		pc.AbortHttp(c, h.BadGateway)
		return
	}

	pc.WriteResponse(c, web.WebResponse{
		Code:    200,
		Message: "success",
	})
}

func (pc *CommentControllerImpl) BulkComment(c *gin.Context) {
	if h.GetStage(c) != "Development" {
		pc.CustomMsgAbortHttp(c, "No Content", 204)
		return
	}

	var datas web.CommentDatas
	pc.GetParams(c, &datas)

	var comments []models.Comment
	var wg sync.WaitGroup
	for _, data := range datas.Datas {
		wg.Add(1)
		go func(data web.CommentData) {
			defer wg.Done()
			postId, _ := primitive.ObjectIDFromHex(data.PostId.Hex())
			t, _ := time.Parse("2006-01-02T15:04:05Z07:00", data.CreatedAt)
			u, _ := time.Parse("2006-01-02T15:04:05Z07:00", data.UpdatedAt)
			comments = append(comments, m.Comment{
				PostId:    postId,
				UserId:    data.UserId,
				CreatedAt: t,
				UpdatedAt: u,
				Text:      h.Encryption(data.Text),
				Reply:     []m.ReplyComment{},
			})
		}(data)
	}

	wg.Wait()
	pc.Service.InsertManyAndBindIds(context.Background(), comments)

	// var commentDocuments []br.CommentDocumment
	// for _, comment := range comments {
	// 	commentDocuments = append(commentDocuments, br.CommentDocumment{
	// 		Id:        comment.Id.Hex(),
	// 		UserId:    comment.UserId,
	// 		PostId:    comment.PostId.Hex(),
	// 		Text:      comment.Text,
	// 		CreatedAt: time.Now(),
	// 		UpdatedAt: time.Now(),
	// 	})
	// }

	// if err := br.Broker.PublishMessage(
	// 	context.Background(),
	// 	br.COMMENTEXCHANGE,
	// 	br.BULKCOMMENTQUEUE,
	// 	"application/json",
	// 	&commentDocuments,
	// ); err != nil {
	// 	web.AbortHttp(c, h.BadGateway)
	// 	return
	// }

	pc.WriteResponse(c, web.WebResponse{
		Code:    201,
		Message: "success",
		Data:    comments,
	})
}

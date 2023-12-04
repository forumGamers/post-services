package controller

import (
	"context"

	"github.com/gin-gonic/gin"
	// br "github.com/post-services/broker"
	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	"github.com/post-services/pkg/comment"
	r "github.com/post-services/pkg/reply"
	"github.com/post-services/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReplyController interface {
	AddReply(c *gin.Context)
	DeleteReply(c *gin.Context)
}

type ReplyControllerImpl struct {
	Service r.ReplyService
	Comment comment.CommentRepo
}

func NewReplyController(
	service r.ReplyService,
	commentRepo comment.CommentRepo,
) ReplyController {
	return &ReplyControllerImpl{
		Service: service,
		Comment: commentRepo,
	}
}

func (rc *ReplyControllerImpl) AddReply(c *gin.Context) {
	commentId, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		web.AbortHttp(c, h.ErrInvalidObjectId)
		return
	}

	var data web.CommentForm
	c.ShouldBind(&data)

	if err := rc.Service.ValidateReply(&data); err != nil {
		web.HttpValidationErr(c, err)
		return
	}

	var comment m.Comment
	if err := rc.Comment.FindById(context.Background(), commentId, &comment); err != nil {
		web.AbortHttp(c, err)
		return
	}

	reply := rc.Service.CreatePayload(data, h.GetUser(c).UUID)
	if err := rc.Comment.CreateReply(context.Background(), comment.Id, &reply); err != nil {
		web.AbortHttp(c, err)
		return
	}

	// if err := br.Broker.PublishMessage(context.Background(), br.REPLYEXCHANGE, br.NEWREPLYQUEUE, "application/json", br.ReplyDocument{
	// 	Id:        reply.Id.Hex(),
	// 	UserId:    reply.UserId,
	// 	Text:      reply.Text,
	// 	CommentId: reply.CommentId.Hex(),
	// 	CreatedAt: reply.CreatedAt,
	// 	UpdatedAt: reply.UpdatedAt,
	// }); err != nil {
	// 	web.AbortHttp(c, h.BadGateway)
	// 	return
	// }

	reply.Text = h.Decryption(reply.Text)

	web.WriteResponse(c, web.WebResponse{
		Code:    201,
		Message: "success",
		Data:    reply,
	})
}

func (rc *ReplyControllerImpl) DeleteReply(c *gin.Context) {
	replyId, err := primitive.ObjectIDFromHex(c.Param("replyId"))
	if err != nil {
		web.AbortHttp(c, h.ErrInvalidObjectId)
		return
	}
	commentId, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		web.AbortHttp(c, h.ErrInvalidObjectId)
		return
	}

	var reply m.ReplyComment
	if err := rc.Comment.FindReplyById(context.Background(), commentId, replyId, &reply); err != nil {
		web.AbortHttp(c, err)
		return
	}

	if err := rc.Service.AuthorizeDeleteReply(reply, h.GetUser(c)); err != nil {
		web.AbortHttp(c, err)
		return
	}

	if err := rc.Comment.DeleteOneReply(context.Background(), commentId, replyId); err != nil {
		web.AbortHttp(c, err)
		return
	}

	// if err := br.Broker.PublishMessage(context.Background(), br.REPLYEXCHANGE, br.DELETEREPLYQUEUE, "application/json", br.ReplyDocument{
	// 	Id:        reply.Id.Hex(),
	// 	UserId:    reply.UserId,
	// 	Text:      reply.Text,
	// 	CommentId: reply.CommentId.Hex(),
	// 	CreatedAt: reply.CreatedAt,
	// 	UpdatedAt: reply.UpdatedAt,
	// }); err != nil {
	// 	web.AbortHttp(c, h.BadGateway)
	// 	return
	// }

	web.WriteResponse(c, web.WebResponse{
		Code:    200,
		Message: "success",
	})
}

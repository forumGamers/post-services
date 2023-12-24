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

func NewReplyController(
	service r.ReplyService,
	commentRepo comment.CommentRepo,
	r web.RequestReader,
	w web.ResponseWriter,
) ReplyController {
	return &ReplyControllerImpl{w, r, service, commentRepo}
}

func (rc *ReplyControllerImpl) AddReply(c *gin.Context) {
	commentId, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		rc.AbortHttp(c, h.ErrInvalidObjectId)
		return
	}

	var data web.CommentForm
	rc.GetParams(c, &data)

	if err := rc.Service.ValidateReply(&data); err != nil {
		rc.HttpValidationErr(c, err)
		return
	}

	var comment m.Comment
	if err := rc.Comment.FindById(context.Background(), commentId, &comment); err != nil {
		rc.AbortHttp(c, err)
		return
	}

	reply := rc.Service.CreatePayload(data, h.GetUser(c).UUID)
	if err := rc.Comment.CreateReply(context.Background(), comment.Id, &reply); err != nil {
		rc.AbortHttp(c, err)
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
	// 	rc.AbortHttp(c, h.BadGateway)
	// 	return
	// }

	reply.Text = h.Decryption(reply.Text)

	rc.WriteResponse(c, web.WebResponse{
		Code:    201,
		Message: "success",
		Data:    reply,
	})
}

func (rc *ReplyControllerImpl) DeleteReply(c *gin.Context) {
	replyId, err := primitive.ObjectIDFromHex(c.Param("replyId"))
	if err != nil {
		rc.AbortHttp(c, h.ErrInvalidObjectId)
		return
	}
	commentId, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		rc.AbortHttp(c, h.ErrInvalidObjectId)
		return
	}

	var reply m.ReplyComment
	if err := rc.Comment.FindReplyById(context.Background(), commentId, replyId, &reply); err != nil {
		rc.AbortHttp(c, err)
		return
	}

	if err := rc.Service.AuthorizeDeleteReply(reply, h.GetUser(c)); err != nil {
		rc.AbortHttp(c, err)
		return
	}

	if err := rc.Comment.DeleteOneReply(context.Background(), commentId, replyId); err != nil {
		rc.AbortHttp(c, err)
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
	// 	rc.AbortHttp(c, h.BadGateway)
	// 	return
	// }

	rc.WriteResponse(c, web.WebResponse{
		Code:    200,
		Message: "success",
	})
}

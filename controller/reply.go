package controller

import (
	"context"

	"github.com/gin-gonic/gin"
	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	b "github.com/post-services/pkg/base"
	r "github.com/post-services/pkg/reply"
	"github.com/post-services/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReplyController interface {
	AddReply(c *gin.Context)
	DeleteReply(c *gin.Context)
}

type ReplyControllerImpl struct {
	Repo    r.ReplyRepo
	Service r.ReplyService
}

func NewReplyController(service r.ReplyService, repo r.ReplyRepo) ReplyController {
	return &ReplyControllerImpl{
		Repo:    repo,
		Service: service,
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
	if err := b.NewBaseRepo(b.GetCollection(b.Comment)).FindOneById(context.Background(), commentId, &comment); err != nil {
		web.AbortHttp(c, err)
		return
	}

	reply := rc.Service.CreatePayload(data, commentId, h.GetUser(c).UUID)
	if err := rc.Repo.CreateReply(context.Background(), &reply); err != nil {
		web.AbortHttp(c, err)
		return
	}

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

	var reply m.ReplyComment
	if err := rc.Repo.FindById(context.Background(), replyId, &reply); err != nil {
		web.AbortHttp(c, err)
		return
	}

	if err := rc.Service.AuthorizeDeleteReply(reply, h.GetUser(c)); err != nil {
		web.AbortHttp(c, err)
		return
	}

	if err := rc.Repo.DeleteOne(context.Background(), replyId); err != nil {
		web.AbortHttp(c, err)
		return
	}

	web.WriteResponse(c, web.WebResponse{
		Code:    200,
		Message: "success",
	})
}

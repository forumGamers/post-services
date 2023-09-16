package controller

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	br "github.com/post-services/broker"
	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	b "github.com/post-services/pkg/base"
	l "github.com/post-services/pkg/like"
	"github.com/post-services/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeController interface {
	LikePost(c *gin.Context)
	UnlikePost(c *gin.Context)
}

type LikeControllerImpl struct {
	Service l.LikeService
	Repo    l.LikeRepo
}

func NewLikeController(service l.LikeService, repo l.LikeRepo) LikeController {
	return &LikeControllerImpl{
		Service: service,
		Repo:    repo,
	}
}

func (lc *LikeControllerImpl) LikePost(c *gin.Context) {
	postId, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		web.AbortHttp(c, h.ErrInvalidObjectId)
		return
	}

	id := h.GetUser(c).UUID
	var post m.Post
	if err := b.NewBaseRepo(b.GetCollection(b.Post)).FindOneById(context.Background(), postId, &post); err != nil {
		web.AbortHttp(c, err)
		return
	}

	var like m.Like
	if err := lc.Repo.GetLikesByUserIdAndPostId(context.Background(), postId, id, &like); err != nil {
		if err != h.NotFount {
			web.AbortHttp(c, err)
			return
		}
	}

	if like.Id != primitive.NilObjectID {
		web.AbortHttp(c, h.Conflict)
		return
	}

	newLike := m.Like{
		UserId:    id,
		PostId:    post.Id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	result, err := lc.Repo.AddLikes(context.Background(), &newLike)
	if err != nil {
		web.AbortHttp(c, err)
		return
	}

	newLike.Id = result
	if err := br.Broker.PublishMessage(context.Background(), br.LIKEEXCHANGE, br.NEWLIKEQUEUE, "application/json", br.LikeDocument{
		Id:        newLike.Id.Hex(),
		UserId:    newLike.UserId,
		PostId:    newLike.PostId.Hex(),
		CreatedAt: newLike.CreatedAt,
		UpdatedAt: newLike.UpdatedAt,
	}); err != nil {
		web.AbortHttp(c, h.BadGateway)
		return
	}

	web.WriteResponse(c, web.WebResponse{
		Code:    201,
		Message: "success",
		Data:    newLike,
	})
}

func (lc *LikeControllerImpl) UnlikePost(c *gin.Context) {
	postId, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		web.AbortHttp(c, h.ErrInvalidObjectId)
		return
	}

	userId := h.GetUser(c).UUID
	var like m.Like

	if err := lc.Repo.GetLikesByUserIdAndPostId(context.Background(), postId, userId, &like); err != nil {
		web.AbortHttp(c, err)
		return
	}

	if err := lc.Repo.DeleteLike(context.Background(), postId, userId); err != nil {
		web.AbortHttp(c, err)
		return
	}

	if err := br.Broker.PublishMessage(context.Background(), br.LIKEEXCHANGE, br.DELETELIKEQUEUE, "application/json", br.LikeDocument{
		Id:        like.Id.Hex(),
		UserId:    like.UserId,
		PostId:    like.PostId.Hex(),
		CreatedAt: like.CreatedAt,
		UpdatedAt: like.UpdatedAt,
	}); err != nil {
		web.AbortHttp(c, h.BadGateway)
		return
	}

	web.WriteResponse(c, web.WebResponse{
		Code:    200,
		Message: "success",
	})
}

package controller

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	h "github.com/post-services/helper"
	"github.com/post-services/models"
	m "github.com/post-services/models"
	b "github.com/post-services/pkg/base"
	l "github.com/post-services/pkg/like"
	"github.com/post-services/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewLikeController(service l.LikeService, repo l.LikeRepo, r web.RequestReader, w web.ResponseWriter) LikeController {
	return &LikeControllerImpl{w, r, service, repo}
}

func (lc *LikeControllerImpl) LikePost(c *gin.Context) {
	postId, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		lc.AbortHttp(c, h.ErrInvalidObjectId)
		return
	}

	id := h.GetUser(c).UUID
	var post m.Post
	if err := b.NewBaseRepo(b.GetCollection(b.Post)).FindOneById(context.Background(), postId, &post); err != nil {
		lc.AbortHttp(c, err)
		return
	}

	var like m.Like
	if err := lc.Repo.GetLikesByUserIdAndPostId(context.Background(), postId, id, &like); err != nil {
		if err != h.NotFount {
			lc.AbortHttp(c, err)
			return
		}
	}

	if like.Id != primitive.NilObjectID {
		lc.AbortHttp(c, h.Conflict)
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
		lc.AbortHttp(c, err)
		return
	}

	newLike.Id = result

	lc.WriteResponse(c, web.WebResponse{
		Code:    201,
		Message: "success",
		Data:    newLike,
	})
}

func (lc *LikeControllerImpl) UnlikePost(c *gin.Context) {
	postId, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		lc.AbortHttp(c, h.ErrInvalidObjectId)
		return
	}

	userId := h.GetUser(c).UUID
	var like m.Like

	if err := lc.Repo.GetLikesByUserIdAndPostId(context.Background(), postId, userId, &like); err != nil {
		lc.AbortHttp(c, err)
		return
	}

	if err := lc.Repo.DeleteLike(context.Background(), postId, userId); err != nil {
		lc.AbortHttp(c, err)
		return
	}

	lc.WriteResponse(c, web.WebResponse{
		Code:    200,
		Message: "success",
	})
}

func (lc *LikeControllerImpl) BulkLikes(c *gin.Context) {
	if h.GetStage(c) != "Development" {
		lc.CustomMsgAbortHttp(c, "No Content", 204)
		return
	}

	var datas web.LikeDatas
	lc.GetParams(c, &datas)

	var likes []models.Like
	var wg sync.WaitGroup
	for _, like := range datas.Datas {
		wg.Add(1)
		go func(like web.LikeData) {
			defer wg.Done()
			postId, _ := primitive.ObjectIDFromHex(like.PostId.Hex())
			t, _ := time.Parse("2006-01-02T15:04:05Z07:00", like.CreatedAt)
			u, _ := time.Parse("2006-01-02T15:04:05Z07:00", like.UpdatedAt)
			likes = append(likes, models.Like{
				PostId:    postId,
				UserId:    like.UserId,
				CreatedAt: t,
				UpdatedAt: u,
			})
		}(like)
	}

	wg.Wait()
	lc.Service.InsertManyAndBindIds(context.Background(), likes)

	lc.WriteResponse(c, web.WebResponse{
		Code:    201,
		Message: "Success",
		Data:    likes,
	})
}

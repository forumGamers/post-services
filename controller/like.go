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

type LikeController interface {
	LikePost(c *gin.Context)
	UnlikePost(c *gin.Context)
	BulkLikes(c *gin.Context)
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

	web.WriteResponse(c, web.WebResponse{
		Code:    200,
		Message: "success",
	})
}

func (lc *LikeControllerImpl) BulkLikes(c *gin.Context) {
	if h.GetStage(c) != "Development" {
		web.CustomMsgAbortHttp(c, "No Content", 204)
		return
	}

	var datas web.LikeDatas
	c.ShouldBind(&datas)

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

	web.WriteResponse(c, web.WebResponse{
		Code:    201,
		Message: "Success",
		Data:    likes,
	})
}

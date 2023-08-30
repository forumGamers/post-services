package controller

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	b "github.com/post-services/pkg/base"
	l "github.com/post-services/pkg/like"
	"github.com/post-services/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeController interface{
	LikePost(c *gin.Context)
}

type LikeControllerImpl struct {
	Service l.LikeService
	Repo 	l.LikeRepo
}

func NewLikeController(service l.LikeService,repo l.LikeRepo) LikeController {
	return &LikeControllerImpl{
		Service: service,
		Repo: repo,
	}
}

func (lc *LikeControllerImpl) LikePost(c *gin.Context) {
	postId,err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		web.AbortHttp(c,h.ErrInvalidObjectId)
		return
	}

	var wg sync.WaitGroup
	errCh := make(chan error)
	id := h.GetUser(c).Id

	wg.Add(2)
	go func ()  {
		defer wg.Done()
		var post m.Post
		errCh <- b.NewBaseRepo(b.GetCollection(b.Post)).FindOneById(context.Background(),postId,&post)
	}()

	go func ()  {
		defer wg.Done()
		var like m.Like
		if err := lc.Repo.GetLikesByUserIdAndPostId(context.Background(),postId,id,&like) ; err != nil {
			if err == h.NotFount {
				errCh <- nil
				return
			}
			errCh <- err
			return
		}
		errCh <- h.Conflict
	}()

	for i := 0 ; i < 2 ; i++ {
		select {
			case err := <- errCh : 
				if err != nil {
					web.AbortHttp(c,err)
					return
				}
		}
	}
	wg.Wait()

	like := m.Like{
		UserId: id,
		PostId: postId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	result,err := lc.Repo.AddLikes(context.Background(),&like)
	if err != nil {
		web.AbortHttp(c,err)
		return
	}

	like.Id = result
	web.WriteResponse(c,web.WebResponse{
		Code: 201,
		Message: "success",
		Data: like,
	})
}
package controller

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	"github.com/post-services/pkg/bookmark"
	"github.com/post-services/pkg/post"
	"github.com/post-services/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewBookmarkController(
	repo bookmark.BookmarkRepo,
	service bookmark.BookmarkService,
	postRepo post.PostRepo,
	w web.ResponseWriter,
	r web.RequestReader,
) BookmarkController {
	return &BookmarkControllerImpl{w, r, repo, service, postRepo}
}

func (bc *BookmarkControllerImpl) CreateBookmark(c *gin.Context) {
	userId := h.GetUser(c).UUID
	postId, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		bc.AbortHttp(c, bc.NewInvalidObjectIdError())
		return
	}

	var post m.Post
	if err := bc.PostRepo.FindById(context.Background(), postId, &post); err != nil {
		bc.AbortHttp(c, err)
		return
	}

	var bookmark bookmark.Bookmark
	if err := bc.Repo.FindOne(context.Background(), bson.M{"postId": postId, "userId": userId}, &bookmark); err != nil {
		if err != mongo.ErrNoDocuments || strings.ToLower(err.Error()) != "data not found" {
			bc.AbortHttp(c, err)
			return
		}
	} else {
		bc.AbortHttp(c, bc.New409Error("Conflict"))
		return
	}

	data := bc.Service.CreatePayload(postId, userId)
	if err := bc.Repo.CreateOne(context.Background(), &data); err != nil {
		bc.AbortHttp(c, err)
		return
	}

	bc.Write200Response(c, "success", data)
}

package controller

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	r "github.com/post-services/repository"
	s "github.com/post-services/services"
	tp "github.com/post-services/third-party"
	"github.com/post-services/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostController interface {
	CreatePost(c *gin.Context)
	FindById(c *gin.Context)
}

type PostControllerImpl struct {
	Service		s.PostService
	Repo		r.PostRepo
}

func NewPostController(service s.PostService,repo r.PostRepo) PostController {
	return &PostControllerImpl{
		Service: service,
		Repo: repo,
	}
}

func (pc *PostControllerImpl) CreatePost(c *gin.Context){
	var form web.PostForm
	c.ShouldBind(&form)

	if err := pc.Service.ValidatePostInput(&form) ; err != nil {
		web.HttpValidationErr(c,err)
		return
	}

	user := h.GetUser(c)
	var media []byte = nil
	var folderName string
	fileName := ""
	// var savedFile *os.File

	if form.File != nil {
		var err error = nil
		media,_ ,err = h.SaveUploadedFile(c,form.File)
		if err != nil {
			panic(err.Error())
		}

		folderName,err = h.CheckFileType(form.File)
		if err != nil {
			panic(err.Error())
		}

		fileName = form.File.Filename
	}

	data,err := pc.Service.CreatePostPayload(context.Background(),&form,user,tp.UploadFile{
		Data: media,
		Folder: folderName,
		Name: fileName,
	})

	if err != nil {
		web.AbortHttp(c,err)
	}

	pc.Repo.Create(context.Background(),&data)

	if form.File != nil {
		// savedFile.Close()
		os.Remove(h.GetUploadDir(fileName))
	}

	web.WriteResponse(c,web.WebResponse{
		Code: 201,
		Message: "Success",
		Data: data,
	})
}

func (pc *PostControllerImpl) FindById(c *gin.Context){
	postId,err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		web.AbortHttp(c,h.ErrInvalidObjectId)
		return
	}

	var data m.Post
	if err := pc.Repo.FindById(context.Background(),postId,&data) ; err != nil {
		web.AbortHttp(c,err)
		return
	}

	web.WriteResponse(c,web.WebResponse{
		Data: data,
		Message: "Success",
		Code: 200,
	})
}

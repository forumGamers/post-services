package controller

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	tp "github.com/post-services/third-party"
	p "github.com/post-services/pkg/post"
	"github.com/post-services/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostController interface {
	CreatePost(c *gin.Context)
	FindById(c *gin.Context)
}

type PostControllerImpl struct {
	Service		p.PostService
	Repo		p.PostRepo
}

func NewPostController(service p.PostService,repo p.PostRepo) PostController {
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
	fileInfo := struct{
		Media []byte
		FolderName string
		FileName string
		SavedFile *os.File
	}{}

	if form.File != nil {
		var err error = nil
		fileInfo.Media,fileInfo.SavedFile ,err = h.SaveUploadedFile(c,form.File)
		if err != nil {
			panic(err.Error())
		}

		fileInfo.FolderName,err = h.CheckFileType(form.File)
		if err != nil {
			panic(err.Error())
		}

		fileInfo.FolderName = "post_"+fileInfo.FolderName
		fileInfo.FileName = form.File.Filename
	}

	data,err := pc.Service.CreatePostPayload(context.Background(),&form,user,tp.UploadFile{
		Data: fileInfo.Media,
		Folder: fileInfo.FolderName,
		Name: fileInfo.FileName,
	})

	if err != nil {
		web.AbortHttp(c,err)
		return
	}

	pc.Repo.Create(context.Background(),&data)

	if form.File != nil {
		fileInfo.SavedFile.Close()
		os.Remove(h.GetUploadDir(fileInfo.FileName))
	}

	data.Text = h.Decryption(data.Text)

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

	data.Text = h.Decryption(data.Text)

	web.WriteResponse(c,web.WebResponse{
		Data: data,
		Message: "Success",
		Code: 200,
	})
}

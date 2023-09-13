package controller

import (
	"context"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	br "github.com/post-services/broker"
	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	b "github.com/post-services/pkg/base"
	p "github.com/post-services/pkg/post"
	"github.com/post-services/pkg/reply"
	tp "github.com/post-services/third-party"
	"github.com/post-services/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostController interface {
	CreatePost(c *gin.Context)
	FindById(c *gin.Context)
	DeletePost(c *gin.Context)
}

type PostControllerImpl struct {
	Service p.PostService
	Repo    p.PostRepo
}

func NewPostController(service p.PostService, repo p.PostRepo) PostController {
	return &PostControllerImpl{
		Service: service,
		Repo:    repo,
	}
}

func (pc *PostControllerImpl) CreatePost(c *gin.Context) {
	var form web.PostForm
	c.ShouldBind(&form)

	if err := pc.Service.ValidatePostInput(&form); err != nil {
		web.HttpValidationErr(c, err)
		return
	}

	user := h.GetUser(c)
	fileInfo := struct {
		Media      []byte
		FolderName string
		FileName   string
		SavedFile  *os.File
	}{}

	if form.File != nil {
		var err error = nil
		fileInfo.Media, fileInfo.SavedFile, err = h.SaveUploadedFile(c, form.File)
		if err != nil {
			panic(err.Error())
		}

		fileInfo.FolderName, err = h.CheckFileType(form.File)
		if err != nil {
			panic(err.Error())
		}

		fileInfo.FolderName = "post_" + fileInfo.FolderName
		fileInfo.FileName = form.File.Filename
	}

	data, err := pc.Service.CreatePostPayload(context.Background(), &form, user, tp.UploadFile{
		Data:   fileInfo.Media,
		Folder: fileInfo.FolderName,
		Name:   fileInfo.FileName,
	})

	if err != nil {
		web.AbortHttp(c, err)
		return
	}

	pc.Repo.Create(context.Background(), &data)

	if form.File != nil {
		fileInfo.SavedFile.Close()
		os.Remove(h.GetUploadDir(fileInfo.FileName))
	}

	if err := br.Broker.PublishMessage(context.Background(), br.POSTEXCHANGE, br.NEWPOSTQUEUE, "application/json", br.PostDocument{
		Id:           data.Id.Hex(),
		UserId:       data.UserId,
		Text:         data.Text,
		AllowComment: data.AllowComment,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
		Tags:         data.Tags,
		Privacy:      data.Privacy,
		Media:        br.Media(data.Media),
	}); err != nil {
		//handle koneksi nya putus
		web.AbortHttp(c, h.InternalServer)
		return
	}

	data.Text = h.Decryption(data.Text)

	web.WriteResponse(c, web.WebResponse{
		Code:    201,
		Message: "Success",
		Data:    data,
	})
}

func (pc *PostControllerImpl) DeletePost(c *gin.Context) {
	postId, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		web.AbortHttp(c, h.ErrInvalidObjectId)
		return
	}

	var data m.Post
	if err := pc.Repo.FindById(context.Background(), postId, &data); err != nil {
		web.AbortHttp(c, err)
		return
	}

	user := h.GetUser(c)

	if data.UserId != user.UUID || user.Role != "Admin" {
		web.AbortHttp(c, h.Forbidden)
		return
	}

	session, err := pc.Repo.GetSession()
	if err != nil {
		web.AbortHttp(c, err)
		return
	}

	defer session.EndSession(context.Background())

	if err := session.StartTransaction(); err != nil {
		web.AbortHttp(c, err)
		return
	}

	ctx := mongo.NewSessionContext(context.Background(), session)
	var wg sync.WaitGroup
	errCh := make(chan error)
	filter := bson.M{"postId": data.Id}
	wg.Add(6)
	runRountine := func(f func()) {
		defer wg.Done()
		f()
	}

	go runRountine(func() {
		pc.Service.DeletePostMedia(ctx, data, errCh)
	})
	go runRountine(func() {
		errCh <- b.NewBaseRepo(b.GetCollection(b.Like)).DeleteMany(ctx, filter)
	})
	go runRountine(func() {
		errCh <- b.NewBaseRepo(b.GetCollection(b.Comment)).DeleteMany(ctx, filter)
	})
	go runRountine(func() {
		errCh <- pc.Repo.DeleteOne(ctx, data.Id)
	})
	go runRountine(func() {
		errCh <- b.NewBaseRepo(b.GetCollection(b.Share)).DeleteMany(ctx, filter)
	})
	go runRountine(func() {
		errCh <- reply.NewReplyRepo().DeleteReplyByPostId(ctx, data.Id)
	})

	flag := false
	var errDb error
	for i := 0; i < 6; i++ {
		select {
		case err := <-errCh:
			{
				if err != nil && !flag {
					flag = true
					errDb = err
				}
			}
		}
	}

	wg.Wait()
	if flag {
		session.AbortTransaction(ctx)
		web.AbortHttp(c, errDb)
		return
	}

	if err := br.Broker.PublishMessage(ctx, br.POSTEXCHANGE, br.DELETEPOSTQUEUE, "application/json", br.PostDocument{
		Id:           data.Id.Hex(),
		UserId:       data.UserId,
		Text:         data.Text,
		AllowComment: data.AllowComment,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
		Tags:         data.Tags,
		Privacy:      data.Privacy,
		Media:        br.Media(data.Media),
	}); err != nil {
		println(err.Error(), "channel")
		session.AbortTransaction(ctx)
		web.AbortHttp(c, h.InternalServer)
		return
	}

	if err := session.CommitTransaction(ctx); err != nil {
		println(err.Error(), "commit")
		session.AbortTransaction(ctx)
		web.AbortHttp(c, err)
		return
	}

	web.WriteResponse(c, web.WebResponse{
		Message: "success",
		Code:    200,
	})
}

func (pc *PostControllerImpl) FindById(c *gin.Context) {
	postId, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		web.AbortHttp(c, h.ErrInvalidObjectId)
		return
	}

	var data m.Post
	if err := pc.Repo.FindById(context.Background(), postId, &data); err != nil {
		web.AbortHttp(c, err)
		return
	}

	data.Text = h.Decryption(data.Text)

	web.WriteResponse(c, web.WebResponse{
		Data:    data,
		Message: "Success",
		Code:    200,
	})
}

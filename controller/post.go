package controller

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	// br "github.com/post-services/broker"
	h "github.com/post-services/helper"
	"github.com/post-services/models"
	m "github.com/post-services/models"
	"github.com/post-services/pkg/comment"
	"github.com/post-services/pkg/like"
	p "github.com/post-services/pkg/post"
	"github.com/post-services/pkg/share"
	tp "github.com/post-services/third-party"
	"github.com/post-services/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewPostController(
	service p.PostService,
	repo p.PostRepo,
	commentRepo comment.CommentRepo,
	likeRepo like.LikeRepo,
	shareRepo share.ShareRepo,
	r web.RequestReader,
	w web.ResponseWriter,
) PostController {
	return &PostControllerImpl{
		w, r, service, repo, commentRepo, likeRepo, shareRepo,
	}
}

func (pc *PostControllerImpl) CreatePost(c *gin.Context) {
	var form web.PostForm
	pc.GetParams(c, &form)

	if err := pc.Service.ValidatePostInput(&form); err != nil {
		pc.HttpValidationErr(c, err)
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
			pc.New400Error(err.Error())
		}

		fileInfo.FolderName, err = h.CheckFileType(form.File)
		if err != nil {
			pc.New400Error(err.Error())
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
		pc.AbortHttp(c, err)
		return
	}

	pc.Repo.Create(context.Background(), &data)

	if form.File != nil {
		fileInfo.SavedFile.Close()
		os.Remove(h.GetUploadDir(fileInfo.FileName))
	}

	// go br.Broker.PublishMessage(context.Background(), br.POSTEXCHANGE, br.NEWPOSTQUEUE, "application/json", br.PostDocument{
	// 	Id:           data.Id.Hex(),
	// 	UserId:       data.UserId,
	// 	Text:         data.Text,
	// 	AllowComment: data.AllowComment,
	// 	CreatedAt:    data.CreatedAt,
	// 	UpdatedAt:    data.UpdatedAt,
	// 	Tags:         data.Tags,
	// 	Privacy:      data.Privacy,
	// 	Media:        br.Media(data.Media),
	// })

	data.Text = h.Decryption(data.Text)

	pc.WriteResponse(c, web.WebResponse{
		Code:    201,
		Message: "Success",
		Data:    data,
	})
}

func (pc *PostControllerImpl) DeletePost(c *gin.Context) {
	postId, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		pc.AbortHttp(c, pc.NewInvalidObjectIdError())
		return
	}

	var data m.Post
	if err := pc.Repo.FindById(context.Background(), postId, &data); err != nil {
		pc.AbortHttp(c, err)
		return
	}

	user := h.GetUser(c)

	if data.UserId != user.UUID || user.LoggedAs != "Admin" {
		pc.AbortHttp(c, pc.New403Error("Forbidden"))
		return
	}

	session, err := pc.Repo.GetSession()
	if err != nil {
		pc.AbortHttp(c, err)
		return
	}

	defer session.EndSession(context.Background())

	if err := session.StartTransaction(); err != nil {
		pc.AbortHttp(c, err)
		return
	}

	ctx := mongo.NewSessionContext(context.Background(), session)
	var wg sync.WaitGroup
	errCh := make(chan error)
	wg.Add(5)
	runRountine := func(f func()) {
		defer wg.Done()
		f()
	}

	go runRountine(func() {
		pc.Service.DeletePostMedia(ctx, data, errCh)
	})
	go runRountine(func() {
		errCh <- pc.LikeRepo.DeletePostLikes(ctx, data.Id)
	})
	go runRountine(func() {
		errCh <- pc.CommentRepo.DeleteMany(ctx, data.Id)
	})
	go runRountine(func() {
		errCh <- pc.Repo.DeleteOne(ctx, data.Id)
	})
	go runRountine(func() {
		errCh <- pc.ShareRepo.DeleteMany(ctx, data.Id)
	})

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			session.AbortTransaction(ctx)
			pc.AbortHttp(c, err)
			return
		}
	}

	// go br.Broker.PublishMessage(ctx, br.POSTEXCHANGE, br.DELETEPOSTQUEUE, "application/json", br.PostDocument{
	// 	Id:           data.Id.Hex(),
	// 	UserId:       data.UserId,
	// 	Text:         data.Text,
	// 	AllowComment: data.AllowComment,
	// 	CreatedAt:    data.CreatedAt,
	// 	UpdatedAt:    data.UpdatedAt,
	// 	Tags:         data.Tags,
	// 	Privacy:      data.Privacy,
	// 	Media:        br.Media(data.Media),
	// })

	if err := session.CommitTransaction(ctx); err != nil {
		session.AbortTransaction(ctx)
		pc.AbortHttp(c, err)
		return
	}

	pc.WriteResponse(c, web.WebResponse{
		Message: "success",
		Code:    200,
	})
}

func (pc *PostControllerImpl) BulkCreatePost(c *gin.Context) {
	if h.GetStage(c) != "Development" {
		pc.CustomMsgAbortHttp(c, "No Content", 204)
		return
	}

	var datas web.PostDatas
	c.ShouldBind(&datas)

	var posts []models.Post
	var wg sync.WaitGroup
	for _, data := range datas.Datas {
		wg.Add(1)
		go func(data web.PostData) {
			defer wg.Done()
			t, _ := time.Parse("2006-01-02T15:04:05Z07:00", data.CreatedAt)
			u, _ := time.Parse("2006-01-02T15:04:05Z07:00", data.UpdatedAt)
			data.Text = h.Encryption(data.Text)
			posts = append(posts, models.Post{
				UserId: data.UserId,
				Text:   data.Text,
				Media: models.Media{
					Url:  data.Media.Url,
					Id:   data.Media.Id,
					Type: data.Media.Type,
				},
				AllowComment: data.AllowComment,
				Tags:         []string{},
				Privacy:      data.Privacy,
				CreatedAt:    t,
				UpdatedAt:    u,
			})
		}(data)
	}
	wg.Wait()

	pc.Service.InsertManyAndBindIds(context.Background(), posts)

	// go br.Broker.PublishMessage(
	// 	context.Background(),
	// 	br.POSTEXCHANGE,
	// 	br.BULKPOSTQUEUE,
	// 	"application/json",
	// 	&postDocuments,
	// )

	pc.WriteResponse(c, web.WebResponse{
		Code:    201,
		Message: "Success",
		Data:    posts,
	})
}

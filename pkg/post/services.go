package post

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	tp "github.com/post-services/third-party"
	"github.com/post-services/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService interface {
	ValidatePostInput(data *web.PostForm) error
	CreatePostPayload(ctx context.Context, data *web.PostForm, user m.User, file tp.UploadFile) (m.Post, error)
	DeletePostMedia(ctx context.Context, post m.Post, ch chan<- error)
	InsertManyAndBindIds(ctx context.Context, datas []web.PostData) error
}

type PostServiceImpl struct {
	Repo     PostRepo
	Validate *validator.Validate
	ImageKit tp.ImageKitService
}

func NewPostService(repo PostRepo, validate *validator.Validate, ik tp.ImageKitService) PostService {
	return &PostServiceImpl{
		Repo:     repo,
		Validate: validate,
		ImageKit: ik,
	}
}

func (ps *PostServiceImpl) ValidatePostInput(data *web.PostForm) error {
	return ps.Validate.Struct(data)
}

func (ps *PostServiceImpl) CreatePostPayload(
	ctx context.Context,
	data *web.PostForm,
	user m.User,
	file tp.UploadFile,
) (m.Post, error) {
	imageKitCh := make(chan tp.ImageKitResult)

	if file.Data != nil {
		go func() {
			imageKitCh <- ps.ImageKit.UploadFile(ctx, file.Data, file.Name, file.Folder)
		}()
	} else {
		go func() {
			imageKitCh <- tp.ImageKitResult{Url: "", FileId: "", Error: nil}
		}()
	}

	result := <-imageKitCh

	if result.Error != nil {
		return m.Post{}, h.BadGateway
	}

	return m.Post{
		UserId: user.UUID,
		Text:   h.Encryption(data.Text),
		Media: m.Media{
			Url:  result.Url,
			Id:   result.FileId,
			Type: file.Folder,
		},
		AllowComment: data.AllowComment,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Privacy:      data.Privacy,
	}, nil
}

func (ps *PostServiceImpl) DeletePostMedia(ctx context.Context, post m.Post, ch chan<- error) {
	if post.Media.Id != "" {
		go ps.ImageKit.Delete(ctx, post.Media.Id, ch)
	} else {
		go func() {
			ch <- nil
		}()
	}
}

func (ps *PostServiceImpl) InsertManyAndBindIds(ctx context.Context, datas []web.PostData) error {
	var payload []any

	for _, data := range datas {
		payload = append(payload, data)
	}

	ids, err := ps.Repo.CreateMany(ctx, payload)
	if err != nil {
		return err
	}

	for i := 0; i < len(ids.InsertedIDs); i++ {
		id := ids.InsertedIDs[i].(primitive.ObjectID)
		datas[i].Id = id
	}
	return nil
}

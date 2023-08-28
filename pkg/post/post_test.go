package post_test

import (
	"mime/multipart"
	"testing"

	p "github.com/post-services/pkg/post"
	post_utils "github.com/post-services/pkg/post/utils"
	v "github.com/post-services/validations"
	"github.com/post-services/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockRepo = p.PostRepoMockImpl{Mock: mock.Mock{}}
var mockImageKit = p.ImagekitMockImpl{Mock: mock.Mock{}}

var postTest = &p.PostServiceImpl{
	Repo: &mockRepo,
	Validate: v.GetValidator(),
	ImageKit: &mockImageKit,
}

func TestMain(m *testing.M){
	m.Run()
}

func TestValidatePostInput(t *testing.T) {
	file,_ := post_utils.ReadFile("../static/meteor2-02.png")

	datas := []struct{
		Name		string
		Data		web.PostForm
		Error		bool
	}{
		{
			Name: "File = nil & Text = ''",
			Data: web.PostForm{
				File: nil,
				Text: "",
				Privacy: "Public",
				AllowComment: true,
			},
			Error: true,
		},
		{
			Name: "Privacy input isnot in enum",
			Data: web.PostForm{
				File: nil,
				Text: "test",
				Privacy: "Close",
				AllowComment: true,
			},
			Error: true,
		},
		{
			Name: "success with file",
			Data: web.PostForm{
				File: &multipart.FileHeader{
					Filename: "image",
					Size: int64(len(file)),
				},
				Text: "test",
				Privacy: "Public",
				AllowComment: true,
			},
			Error: false,
		},
		{
			Name: "success with file only",
			Data: web.PostForm{
				File: &multipart.FileHeader{
					Filename: "image",
					Size: int64(len(file)),
				},
				Text: "",
				Privacy: "Public",
				AllowComment: true,
			},
			Error: false,
		},
		{
			Name: "success with text only",
			Data: web.PostForm{
				File: nil,
				Text: "test",
				Privacy: "Public",
				AllowComment: true,
			},
			Error: false,
		},
	}

	for _,data := range datas {
		t.Run(data.Name,func(t *testing.T) {
			err := postTest.ValidatePostInput(&data.Data)

			if data.Error {
				assert.NotNil(t,err)
			}else {
				assert.Nil(t,err)
			}
		})
	}
}

func BenchmarkValidatePostInput(b *testing.B) {
	file,_ := post_utils.ReadFile("../static/meteor2-02.png")

	data := web.PostForm{
		File: &multipart.FileHeader{
			Filename: "image",
			Size: int64(len(file)),
		},
		Text: "test",
		Privacy: "Public",
		AllowComment: true,
	}

	b.ResetTimer()

	for i := 0 ; i < b.N ; i++ {
		if err := postTest.ValidatePostInput(&data) ; err != nil {
			b.Fatal(err)
		}
	}
}


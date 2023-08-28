package post_test

import (
	"testing"

	p "github.com/post-services/pkg/post"
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

func TestValidatePostInput(t *testing.T) {
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
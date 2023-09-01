package middlewares

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	b "github.com/post-services/pkg/base"
)

func Logging(c *gin.Context) {
	defer func() {
		if _, err := b.NewBaseRepo(b.GetCollection(b.Log)).Create(context.Background(), m.Log{
			Path:         c.Request.URL.Path,
			UserId:       h.GetUser(c).Id,
			Method:       c.Request.Method,
			StatusCode:   c.Writer.Status(),
			Origin:       c.Request.Header.Get("Origin"),
			ResponseTime: int(time.Since(c.MustGet("start").(time.Time)).Milliseconds()),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}); err != nil {
			println(err)
			return
		}
	}()
	c.Next()
}

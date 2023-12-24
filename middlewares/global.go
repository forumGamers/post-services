package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/post-services/web"
)

func NewMiddlewares(w web.ResponseWriter, r web.RequestReader) Middleware {
	return &MiddlewareImpl{w, r}
}

func (m *MiddlewareImpl) Next(c *gin.Context) {
	c.Next()
}

func (m *MiddlewareImpl) SetContext(c *gin.Context, key string, value any) {
	c.Set(key, value)
}

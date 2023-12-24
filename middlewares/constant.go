package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/post-services/web"
)

type Middleware interface {
	Authentication(c *gin.Context)
	Cors() gin.HandlerFunc
	SetContexts(c *gin.Context)
	CheckOrigin(c *gin.Context)
}

type MiddlewareImpl struct {
	web.ResponseWriter
	web.RequestReader
}

package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
)

func SetStart(c *gin.Context) {
	c.Set("start", time.Now())

	c.Next()
}

package middlewares

import (
	"os"

	"github.com/gin-gonic/gin"
)

func SetContext(c *gin.Context) {
	stage := os.Getenv("APP_STAGE")
	if stage == "" {
		stage = "Development"
	}

	c.Set("stage", stage)
	c.Next()
}

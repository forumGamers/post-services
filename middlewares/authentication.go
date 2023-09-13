package middlewares

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	h "github.com/post-services/helper"
	"github.com/post-services/web"
)

func Authentication(c *gin.Context) {
	access_token := c.Request.Header.Get("access_token")
	if access_token == "" {
		web.AbortHttp(c, h.Forbidden)
		return
	}

	claim := jwt.MapClaims{}

	if token, err := jwt.ParseWithClaims(access_token, &claim, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	}); err != nil || !token.Valid {
		web.AbortHttp(c, h.InvalidToken)
		return
	}

	c.Set("user", claim)
	c.Next()
}

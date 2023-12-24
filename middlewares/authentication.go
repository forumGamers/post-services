package middlewares

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	h "github.com/post-services/helper"
)

func (m *MiddlewareImpl) Authentication(c *gin.Context) {
	access_token := c.Request.Header.Get("access_token")
	if access_token == "" {
		m.AbortHttp(c, m.New403Error("Forbidden"))
		return
	}

	claim := jwt.MapClaims{}

	if token, err := jwt.ParseWithClaims(access_token, &claim, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	}); err != nil || !token.Valid {
		m.AbortHttp(c, h.InvalidToken)
		return
	}

	m.SetContext(c, "user", claim)
	m.Next(c)
}

package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	m "github.com/post-services/models"
)

func GetUser(c *gin.Context) m.User {
	var user m.User

	claimMap, ok := c.Get("user")
	if !ok {
		return user
	}

	claim, oke := claimMap.(jwt.MapClaims)
	if !oke {
		return user
	}

	for key, val := range claim {
		switch key {
		case "UUID":
			user.UUID = val.(string)
		case "loggedAs":
			user.LoggedAs = val.(string)
		}
	}
	return user
}

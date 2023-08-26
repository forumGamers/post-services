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
		case "email":
			user.Email = val.(string)
		case "fullName":
			user.Fullname = val.(string)
		case "iat":
			user.Iat = int(val.(float64))
		case "id":
			user.Id = int(val.(float64))
		case "isVerified":
			user.IsVerified = val.(bool)
		case "phoneNumber":
			user.PhoneNumber = val.(string)
		case "username":
			user.Username = val.(string)
		case "StoreId":
			user.StoreId = int(val.(float64))
		case "role":
			user.Role = val.(string)
		case "point":
			user.Point = int(val.(float64))
		case "experience":
			user.Exp = int(val.(float64))
		}
	}
	return user
}
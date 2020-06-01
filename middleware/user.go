package middleware

import (
	"encoding/base64"
	"fmt"
	"gallery0api/models"
	"hash"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireUser(us models.UserService, hmac hash.Hash) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := GetToken(c)
		if token == "" {
			c.Status(401)
			c.Abort()
			return
		}
		hmac.Write([]byte(token))
		hash := hmac.Sum(nil)
		hmac.Reset()
		fmt.Println("hash", base64.URLEncoding.EncodeToString(hash))
		encode := base64.URLEncoding.EncodeToString(hash)
		user, err := us.GetByToken(encode)
		if err != nil {
			c.Status(401)
			c.Abort()
			return
		}
		c.Set("user", user)
	}
}
func GetToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	header = strings.TrimSpace(header)
	min := len("Bearer ")
	if len(header) <= min {
		return ""
	}
	return header[min:]
}

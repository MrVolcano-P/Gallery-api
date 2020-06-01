package middleware

import (
	"gallery0api/context"
	"gallery0api/header"
	"gallery0api/models"

	"github.com/gin-gonic/gin"
)

func RequireUser(us models.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := header.GetToken(c)
		if token == "" {
			c.Status(401)
			c.Abort()
			return
		}
		user, err := us.GetByToken(token)
		if err != nil {
			c.Status(401)
			c.Abort()
			return
		}
		context.SetUser(c, user)
	}
}

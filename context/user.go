package context

import (
	"gallery0api/models"

	"github.com/gin-gonic/gin"
)

func User(c *gin.Context) *models.UserTable {
	user, ok := c.Value("user").(*models.UserTable)
	if !ok {
		return nil
	}
	return user
}

package middleware

import (
	"gallery0api/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// func Auth(c *gin.Context) {
// 	header := c.GetHeader("Authorization")
// 	checkToken := strings.Split(header, " ")
// 	if len(checkToken) <= 1 {
// 		c.JSON(401, gin.H{
// 			"message": "no token",
// 		})
// 		c.Abort()
// 		return
// 	}
// 	fmt.Println("kkkkkkkkk: ")
// 	token := header[8:]
// 	fmt.Println("token123: ", token)
// 	c.Set("token", token)

// }
func RequireUser(us models.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := GetToken(c)
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

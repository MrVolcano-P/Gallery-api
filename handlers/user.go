package handlers

import (
	"gallery0api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UserHandler struct {
	ug models.UserService
}

func NewUserHandler(ug models.UserService) *UserHandler {
	return &UserHandler{ug}
}

type NewUser struct {
	Email    string
	Password string
	Name     string
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
	cost := 10
	user := &NewUser{}
	if err := c.BindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	userTable := &models.UserTable{}
	userTable.Email = user.Email
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), cost)
	if err != nil {
		c.Status(500)
		return
	}
	userTable.Password = string(hash)
	userTable.Name = user.Name
	if err := uh.ug.CreateUser(userTable); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusCreated, User{
		ID:    userTable.ID,
		Name:  userTable.Name,
		Email: userTable.Email,
	})
}

type LoginReq struct {
	Email    string
	Password string
}

func (uh *UserHandler) Login(c *gin.Context) {
	req := &LoginReq{}
	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	userTable := &models.UserTable{}
	userTable.Email = req.Email
	userTable.Password = req.Password
	token, err := uh.ug.Login(userTable)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}
	userTable.Token = token
	c.JSON(201, gin.H{
		"token": userTable.Token,
	})

}

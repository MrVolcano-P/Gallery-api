package handlers

import (
	"errors"
	"fmt"
	"gallery0api/context"
	"gallery0api/models"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
type SignupReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (h *Handler) Signup(c *gin.Context) {
	req := new(SignupReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	user := new(models.User)
	user.Email = req.Email
	user.Password = req.Password
	user.Name = req.Name
	if err := h.us.Create(user); err != nil {
		Error(c, 500, err)
		return
	}
	c.JSON(201, gin.H{
		"token": user.Token,
		"email": user.Email,
		"name":  user.Name,
	})
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(c *gin.Context) {
	req := new(LoginReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	user := new(models.User)
	user.Email = req.Email
	user.Password = req.Password
	token, err := h.us.Login(user)
	if err != nil {
		Error(c, 401, err)
		return
	}
	c.JSON(201, gin.H{
		"token": token,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	user := context.User(c)
	if user == nil {
		Error(c, 401, errors.New("invalid token"))
		return
	}
	err := h.us.Logout(user)
	if err != nil {
		Error(c, 500, err)
		return
	}
	c.Status(204)
}

func (h *Handler) GetProfile(c *gin.Context) {
	user := context.User(c)
	if user == nil {
		Error(c, 401, errors.New("invalid token"))
		return
	}
	fmt.Println(user)
	c.JSON(200, gin.H{
		"id":user.ID,
		"email":user.Email,
		"name":user.Name,
	})
}

package handlers

import (
	"fmt"
	"gallery0api/context"
	"gallery0api/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// type Handler struct {
// 	gs models.GalleryService
// }

// func NewHandler(gs models.GalleryService) *Handler {
// 	return &Handler{gs}
// }

type GalleryRes struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	IsPublish bool      `json:"is_publish"`
	Owner     User      `json:"owner"`
	CreateAt  time.Time `json:"createAt"`
	UpdateAt  time.Time `json:"updateAt"`
}

type CreateReq struct {
	Name string `json:"name"`
}

type CreateRes struct {
	GalleryRes
}

func (h *Handler) Create(c *gin.Context) {
	user := context.User(c)
	if user == nil {
		c.Status(401)
		return
	}
	req := new(CreateReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	gallery := new(models.Gallery)
	gallery.Name = req.Name
	gallery.UserID = user.ID
	if err := h.gs.Create(gallery); err != nil {
		Error(c, 500, err)
		return
	}
	res := new(CreateRes)
	res.ID = gallery.ID
	res.Name = gallery.Name
	res.IsPublish = gallery.IsPublish
	c.JSON(201, res)
}

func (h *Handler) ListPublish(c *gin.Context) {
	// user := context.User(c)
	// if user == nil {
	// 	c.Status(401)
	// 	return
	// }
	data, err := h.gs.ListAllPublish()
	if err != nil {
		Error(c, 500, err)
		return
	}
	galleries := []GalleryRes{}
	for _, d := range data {
		user, err := h.us.GetByID(d.UserID)
		if err != nil {
			Error(c, 500, err)
			return
		}
		galleries = append(galleries, GalleryRes{
			ID:        d.ID,
			Name:      d.Name,
			IsPublish: d.IsPublish,
			CreateAt:  d.CreatedAt,
			UpdateAt:  d.UpdatedAt,
			Owner: User{
				ID:    user.ID,
				Email: user.Email,
				Name:  user.Name,
			},
		})
	}
	c.JSON(200, galleries)
}

func (h *Handler) List(c *gin.Context) {
	user := context.User(c)
	if user == nil {
		c.Status(401)
		return
	}
	data, err := h.gs.ListByUserID(user.ID)
	if err != nil {
		Error(c, 500, err)
		return
	}
	galleries := []GalleryRes{}
	for _, d := range data {
		user, err := h.us.GetByID(d.UserID)
		if err != nil {
			Error(c, 500, err)
			return
		}
		galleries = append(galleries, GalleryRes{
			ID:        d.ID,
			Name:      d.Name,
			IsPublish: d.IsPublish,
			CreateAt:  d.CreatedAt,
			UpdateAt:  d.UpdatedAt,
			Owner: User{
				ID:    user.ID,
				Email: user.Email,
				Name:  user.Name,
			},
		})
	}
	c.JSON(200, galleries)
}

func (h *Handler) GetOne(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(c, 400, err)
		return
	}
	data, err := h.gs.GetByID(uint(id))
	if err != nil {
		Error(c, 500, err)
		return
	}
	user, err := h.us.GetByID(data.UserID)
	if err != nil {
		Error(c, 500, err)
		return
	}
	if data.IsPublish == true {
		fmt.Println("true")
		c.JSON(200, GalleryRes{
			ID:        data.ID,
			Name:      data.Name,
			IsPublish: data.IsPublish,
			CreateAt:  data.CreatedAt,
			UpdateAt:  data.UpdatedAt,
			Owner: User{
				ID:    user.ID,
				Email: user.Email,
				Name:  user.Name,
			},
		})
	} else {
		fmt.Println("false")
		c.Status(401)
	}
}
func (h *Handler) GetOneAndCheck(c *gin.Context) {
	userContext := context.User(c)
	if userContext == nil {
		c.Status(401)
		return
	}
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(c, 400, err)
		return
	}
	data, err := h.gs.GetByID(uint(id))
	if err != nil {
		Error(c, 500, err)
		return
	}
	if data.UserID != userContext.ID {
		c.JSON(401, gin.H{
			"status":  false,
			"message": "not owner",
		})
		return
	} else {
		c.JSON(200, GalleryRes{
			ID:        data.ID,
			Name:      data.Name,
			IsPublish: data.IsPublish,
			CreateAt:  data.CreatedAt,
			UpdateAt:  data.UpdatedAt,
			Owner: User{
				ID:    userContext.ID,
				Email: userContext.Email,
				Name:  userContext.Name,
			},
		})
	}
}
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(c, 400, err)
		return
	}
	err = h.gs.DeleteGallery(uint(id))
	if err != nil {
		Error(c, 500, err)
		return
	}
	c.Status(204)
}

type UpdateNameReq struct {
	Name string `json:"name"`
}

func (h *Handler) UpdateName(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(c, 400, err)
		return
	}
	req := new(UpdateNameReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	err = h.gs.UpdateGalleryName(uint(id), req.Name)
	if err != nil {
		Error(c, 500, err)
		return
	}
	c.Status(204)
}

type UpdateStatusReq struct {
	IsPublish bool `json:"is_publish"`
}

func (h *Handler) UpdatePublishing(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(c, 400, err)
		return
	}
	req := new(UpdateStatusReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	err = h.gs.UpdateGalleryPublishing(uint(id), req.IsPublish)
	if err != nil {
		Error(c, 500, err)
		return
	}
	c.Status(204)
}

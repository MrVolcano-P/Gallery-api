package handlers

import (
	"fmt"
	"gallery0api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Gallery struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type GalleryHandler struct {
	gg models.GalleryService
}

func NewGalleryHandler(gg models.GalleryService) *GalleryHandler {
	return &GalleryHandler{gg}
}

func (gh *GalleryHandler) ListGallery(c *gin.Context) {
	tts, err := gh.gg.ListGallery()
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	gallerys := []Gallery{}
	for _, tt := range tts {
		gallerys = append(gallerys, Gallery{
			ID:   tt.ID,
			Name: tt.Name,
		})
	}

	c.JSON(http.StatusOK, gallerys)
}

type NewGallery struct {
	Name string
}

func (gh *GalleryHandler) CreateGallery(c *gin.Context) {
	

	data := new(NewGallery)
	if err := c.BindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	GalleryTable := new(models.GalleryTable)
	GalleryTable.Name = data.Name
	if err := gh.gg.CreateGallery(GalleryTable); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusCreated, Gallery{
		ID:   GalleryTable.ID,
		Name: GalleryTable.Name,
	})
}

type UpdateGallery struct {
	Name string
}

func (gh *GalleryHandler) UpdateGallery(c *gin.Context) {
	idString := c.Param("id")
	fmt.Println(idString)
	idUint, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	update := new(UpdateGallery)
	if err := c.BindJSON(update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	GalleryTable := new(models.GalleryTable)
	GalleryTable.Name = update.Name
	if err := gh.gg.UpdateGallery(idUint, GalleryTable); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, Gallery{
		ID:   GalleryTable.ID,
		Name: GalleryTable.Name,
	})

}
func (gh *GalleryHandler) DeleteGallery(c *gin.Context) {
	idString := c.Param("id")
	// idUint, err := strconv.ParseUint(idString, 10, 64)
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	if err := gh.gg.DeleteGallery(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.Status(http.StatusNoContent)
}

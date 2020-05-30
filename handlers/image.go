package handlers

import (
	"fmt"
	"gallery0api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Image struct {
	ID  uint   `json:"id"`
	Src string `json:"src"`
}

type ImageHandler struct {
	ig models.ImageService
}

func NewImageHandler(ig models.ImageService) *ImageHandler {
	return &ImageHandler{ig}
}

type NewImage struct {
	Src       string
	Width     uint
	Height    uint
	GalleryID uint
}

func (ih *ImageHandler) CreateImage(c *gin.Context) {
	getToken, isToken := c.Get("token")
	if !isToken {
		fmt.Println("err")
	}
	fmt.Println(getToken)
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	data := new(NewImage)
	if err := c.BindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	imageTable := new(models.ImageTable)
	imageTable.Src = data.Src
	imageTable.Width = data.Width
	imageTable.Height = data.Height
	imageTable.GalleryID = uint(id)
	if err := ih.ig.CreateImage(imageTable); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusCreated, Image{
		ID:  imageTable.ID,
		Src: imageTable.Src,
	})
}

func (ih *ImageHandler) GetImagesByGalleryID(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	its, err := ih.ig.GetImagesByGalleryID(uint(id))
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	fmt.Println(its)
	// images := []Image{}
	// for _, it := range its {
	// 	images = append(images, Image{
	// 		ID:  it.ID,
	// 		Src: it.Src,
	// 	})
	// }

	// c.JSON(http.StatusOK, images)
}

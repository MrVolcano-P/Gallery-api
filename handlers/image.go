package handlers

import (
	"fmt"
	"gallery0api/models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Image struct {
	ID       uint   `json:"id"`
	Filename string `json:"filename"`
}

type ImageHandler struct {
	ig models.ImageService
}

func NewImageHandler(ig models.ImageService) *ImageHandler {
	return &ImageHandler{ig}
}

type NewImage struct {
	Filename  string
	Width     uint
	Height    uint
	GalleryID uint
}

func (ih *ImageHandler) CreateImage(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = os.MkdirAll("upload", os.ModePerm)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	images := form.File["image"]

	for _, image := range images {
		filename := filepath.Join("image", image.Filename)
		// fmt.Printf("%+v", filename)
		if err := c.SaveUploadedFile(image, filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		imageTable := new(models.ImageTable)
		imageTable.Filename = filename
		imageTable.Width = 4
		imageTable.Height = 3
		imageTable.GalleryID = uint(id)
		if err := ih.ig.CreateImage(imageTable); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusCreated, gin.H{
			
		})
	}
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

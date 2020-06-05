package handlers

import (
	"fmt"
	"gallery0api/models"
	"net/http"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ImageRes struct {
	ID        uint   `json:"id"`
	GalleryID uint   `json:"gallery_id"`
	Filename  string `json:"filename"`
	Width     uint   `json:"width"`
	Height    uint   `json:"height"`
}

type CreateImageRes struct {
	ImageRes
}

// type Handler struct {
// 	gs  models.GalleryService
// 	ims models.ImageService
// }

// func NewHandler(gs models.GalleryService, ims models.ImageService) *Handler {
// 	return &Handler{gs, ims}
// }

func (h *Handler) CreateImage(c *gin.Context) {
	galleryIDStr := c.Param("id")
	galleryID, err := strconv.Atoi(galleryIDStr)
	if err != nil {
		Error(c, 400, err)
		return
	}

	gallery, err := h.gs.GetByID(uint(galleryID))
	if err != nil {
		Error(c, 400, err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		Error(c, 400, err)
		return
	}

	images, err := h.ims.CreateImages(form.File["photos"], gallery.ID)
	if err != nil {
		Error(c, 500, err)
		return
	}

	res := []CreateImageRes{}
	for _, img := range images {
		r := CreateImageRes{}
		r.ID = img.ID
		r.GalleryID = gallery.ID
		r.Filename = path.Join(models.UploadPath, galleryIDStr, img.Filename)
		res = append(res, r)
	}

	c.JSON(201, res)
}

type ListGalleryImagesRes struct {
	ImageRes
}

func (h *Handler) ListGalleryImages(c *gin.Context) {
	galleryIDStr := c.Param("id")
	id, err := strconv.Atoi(galleryIDStr)
	if err != nil {
		Error(c, 400, err)
		return
	}

	gallery, err := h.gs.GetByID(uint(id))
	if err != nil {
		Error(c, 400, err)
		return
	}
	images, err := h.ims.GetByGalleryID(gallery.ID)
	if err != nil {
		Error(c, http.StatusNotFound, err)
		return
	}
	res := []ListGalleryImagesRes{}
	for _, img := range images {
		r := ListGalleryImagesRes{}
		r.ID = img.ID
		r.GalleryID = gallery.ID
		r.Filename = img.FilePath()
		res = append(res, r)
	}
	c.JSON(http.StatusOK, res)
}

type DeleteReq struct {
	FileNames []string `json:"filenames"`
}

func (h *Handler) DeleteImageInGallary(c *gin.Context) {
	galleryIDStr := c.Param("id")
	id, err := strconv.Atoi(galleryIDStr)
	if err != nil {
		Error(c, 400, err)
		return
	}
	req := new(DeleteReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	fmt.Println(req)
	for _, r := range req.FileNames {
		fmt.Println(r)
		err = h.ims.RemoveImageByFileName(uint(id), r)
		if err != nil {
			Error(c, 500, err)
			return
		}
	}
	c.Status(204)
}

func (h *Handler) DeleteImage(c *gin.Context) {
	imageIDStr := c.Param("id")
	filename := c.Param("filename")
	id, err := strconv.Atoi(imageIDStr)
	if err != nil {
		Error(c, 400, err)
		return
	}
	err = h.ims.RemoveImageByFileName(uint(id), filename)
	if err != nil {
		Error(c, 500, err)
		return
	}
	c.Status(200)
	// if err := h.ims.Delete(uint(id)); err != nil {
	// 	Error(c, 500, err)
	// 	return
	// }
}

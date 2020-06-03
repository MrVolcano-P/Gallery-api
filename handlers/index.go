package handlers

import "gallery0api/models"

type Handler struct {
	gs  models.GalleryService
	us  models.UserService
	ims models.ImageService
}

func NewHandler(gs models.GalleryService, us models.UserService, ims models.ImageService) *Handler {
	return &Handler{gs, us, ims}
}

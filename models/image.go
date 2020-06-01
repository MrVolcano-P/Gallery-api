package models

import (
	"github.com/jinzhu/gorm"
)

type ImageTable struct {
	gorm.Model
	Filename  string `gorm:"not null"`
	Width     uint
	Height    uint
	GalleryID uint `gorm:"not null"`
}

type ImageService interface {
	CreateImage(image *ImageTable) error
	GetImagesByGalleryID(id uint) (*ImageTable, error)
}

var _ ImageService = &ImageGorm{}

type ImageGorm struct {
	db *gorm.DB
}

func NewImageGorm(db *gorm.DB) ImageService {
	return &ImageGorm{db}
}

func (ig *ImageGorm) CreateImage(image *ImageTable) error {
	return ig.db.Create(image).Error
}

func (ig *ImageGorm) GetImagesByGalleryID(id uint) (*ImageTable, error) {
	it := new(ImageTable)
	if err := ig.db.Where("gallery_id = ?", id).Find(it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

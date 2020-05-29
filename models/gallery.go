package models

import (
	"github.com/jinzhu/gorm"
)

type GalleryTable struct {
	gorm.Model
	Name   string
	Images []ImageTable
}

type GalleryService interface {
	ListGallery() ([]GalleryTable, error)
	// GetTaskByID(id uint) (*GalleryTable, error)
	CreateGallery(gallery *GalleryTable) error
	UpdateGallery(id uint64, gallery *GalleryTable) error
	DeleteGallery(id uint) error
	
}

var _ GalleryService = &GalleryGorm{}

type GalleryGorm struct {
	db *gorm.DB
}

func NewGalleryGorm(db *gorm.DB) GalleryService {
	return &GalleryGorm{db}
}

func (gg *GalleryGorm) ListGallery() ([]GalleryTable, error) {
	GalleryTables := []GalleryTable{}
	if err := gg.db.Find(&GalleryTables).Error; err != nil {
		return nil, err
	}
	return GalleryTables, nil
}

// func (tg *GalleryGorm) GetTaskByID(id uint) (*GalleryTable, error) {
// 	tt := new(GalleryTable)
// 	if err := tg.db.First(tt, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return tt, nil
// }

func (gg *GalleryGorm) CreateGallery(gallery *GalleryTable) error {
	return gg.db.Create(gallery).Error
}

func (gg *GalleryGorm) UpdateGallery(id uint64, gallery *GalleryTable) error {
	found := new(GalleryTable)
	if err := gg.db.Where("id = ?", id).First(found).Error; err != nil {
		return err
	}
	return gg.db.Model(found).Update("Name", found.Name).Error
}

func (gg *GalleryGorm) DeleteGallery(id uint) error {
	// gt := new(GalleryTable)
	// if err := gg.db.Where("id = ?", id).First(gt).Error; err != nil {
	// 	return err
	// }
	// return gg.db.Delete(gt).Error
	return gg.db.Where("id = ?", id).Delete(&GalleryTable{}).Error
}


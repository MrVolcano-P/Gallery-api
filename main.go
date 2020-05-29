package main

import (
	"gallery0api/handlers"
	"gallery0api/models"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type CreateGallery struct {
	Name string
}

func main() {
	db, err := gorm.Open(
		"mysql",
		"root:password@tcp(poomdv.c52jeww5mzql.ap-southeast-1.rds.amazonaws.com:3306)/gallery?parseTime=true",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.LogMode(true) // dev only!

	if err := db.AutoMigrate(
		&models.GalleryTable{},
		&models.UserTable{},
		&models.ImageTable{},
	).Error; err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	r.Use(cors.New(config))

	gg := models.NewGalleryGorm(db)

	gh := handlers.NewGalleryHandler(gg)

	r.GET("/gallerys", gh.ListGallery)

	r.POST("/gallerys", gh.CreateGallery)

	r.PUT("/gallerys/:id", gh.UpdateGallery)

	r.DELETE("/gallerys/:id", gh.DeleteGallery)

	ug := models.NewUserGorm(db)

	uh := handlers.NewUserHandler(ug)

	r.POST("/signup", uh.CreateUser)

	r.POST("/login", uh.Login)

	ig := models.NewImageGorm(db)

	ih := handlers.NewImageHandler(ig)

	r.POST("/images/:id", ih.CreateImage)

	r.GET("/images/:id", ih.GetImagesByGalleryID)

	r.Run()
}

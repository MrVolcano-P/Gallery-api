package main

import (
	"gallery0api/handlers"
	"gallery0api/middleware"
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

	ug := models.NewUserGorm(db)

	uh := handlers.NewUserHandler(ug)

	r.POST("/signup", uh.CreateUser)

	r.POST("/login", uh.Login)

	ig := models.NewImageGorm(db)

	ih := handlers.NewImageHandler(ig)

	mg := r.Group("")
	mg.Use(middleware.RequireUser(ug))
	{
		mg.GET("/gallerys", gh.ListGallery)

		mg.POST("/gallerys", gh.CreateGallery)

		mg.PUT("/gallerys/:id", gh.UpdateGallery)

		mg.DELETE("/gallerys/:id", gh.DeleteGallery)

		mg.POST("/images/:id", ih.CreateImage)

		mg.GET("/images/:id", ih.GetImagesByGalleryID)

		mg.GET("/sessions", func(c *gin.Context) {
			user, ok := c.Value("user").(*models.UserTable)
			if !ok {
				c.JSON(401, gin.H{
					"message": "invalid token",
				})
				return
			}
			c.JSON(200, user)
		})
	}

	r.Run()
}

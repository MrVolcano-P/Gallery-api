package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"gallery0api/handlers"
	"gallery0api/middleware"
	"gallery0api/models"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

type CreateGallery struct {
	Name string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := gorm.Open(
		"mysql",
		"root:password@tcp(poomdv.c52jeww5mzql.ap-southeast-1.rds.amazonaws.com:3306)/gallery?parseTime=true",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.LogMode(true) // dev only!

	err = models.AutoMigrate(db)
	if err != nil {
		log.Fatal(err)
	}

	hmac := hmac.New(sha256.New, []byte(os.Getenv("Hmackey")))

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = []string{"Authorization", "Origin", "Content-Type"}
	r.Use(cors.New(config))

	gg := models.NewGalleryGorm(db)
	ug := models.NewUserGorm(db, hmac)
	ig := models.NewImageGorm(db)

	gh := handlers.NewGalleryHandler(gg)
	uh := handlers.NewUserHandler(ug)
	ih := handlers.NewImageHandler(ig)

	r.Static("/upload", "./Image")

	r.POST("/signup", uh.CreateUser)
	r.POST("/login", uh.Login)
	r.GET("/gallerys", gh.ListGallery)

	mg := r.Group("")
	mg.Use(middleware.RequireUser(ug, hmac))
	{

		mg.POST("/gallerys", gh.CreateGallery)
		mg.PUT("/gallerys/:id", gh.UpdateGallery)
		mg.DELETE("/gallerys/:id", gh.DeleteGallery)
		mg.POST("/gallerys/:id/images", ih.CreateImage)
		mg.GET("/gallerys/:id/images", ih.GetImagesByGalleryID)
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

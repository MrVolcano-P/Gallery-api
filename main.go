package main

import (
	"gallery0api/handlers"
	"gallery0api/hash"
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

	// err = models.Reset(db)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	
	hmac := hash.NewHMAC(os.Getenv("hmackey"))
	gs := models.NewGalleryService(db)
	ims := models.NewImageService(db)
	us := models.NewUserService(db, hmac)

	gh := handlers.NewGalleryHandler(gs)
	imh := handlers.NewImageHandler(gs, ims)
	uh := handlers.NewUserHandler(us)

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = []string{"Authorization", "Origin", "Content-Type"}

	r.Use(cors.New(config))

	r.Static("/upload", "./upload")

	r.POST("/signup", uh.Signup)
	r.POST("/login", uh.Login)
	r.GET("/galleries", gh.ListPublish)

	auth := r.Group("/")
	auth.Use(middleware.RequireUser(us))
	{
		auth.POST("/logout", uh.Logout)
		user := auth.Group("/user")
		{
			user.POST("/galleries", gh.Create)
			user.GET("/galleries", gh.List)
			user.GET("/galleries/:id", gh.GetOne)
			user.DELETE("/galleries/:id", gh.Delete)
			user.PATCH("/galleries/:id/names", gh.UpdateName)
			user.PATCH("/galleries/:id/publishes", gh.UpdatePublishing)
			user.POST("/galleries/:id/images", imh.CreateImage)
			user.GET("/galleries/:id/images", imh.ListGalleryImages)
			user.DELETE("/images/:id", imh.DeleteImage)
		}

	}
	r.Run(":8080")
}

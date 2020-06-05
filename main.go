package main

import (
	"gallery0api/config"
	"gallery0api/handlers"
	"gallery0api/hash"
	"gallery0api/middleware"
	"gallery0api/models"
	"log"
	"os"
	"time"

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
	conf := config.Load()
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

	if conf.Mode == "dev" {
		db.LogMode(true) // dev only!
	}

	// err = models.Reset(db)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	hmac := hash.NewHMAC(os.Getenv("hmackey"))
	gs := models.NewGalleryService(db)
	ims := models.NewImageService(db)
	us := models.NewUserService(db, hmac)

	// gh := handlers.NewGalleryHandler(gs)
	// imh := handlers.NewImageHandler(gs, ims)
	// uh := handlers.NewUserHandler(us)
	h := handlers.NewHandler(gs, us, ims)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "PUT", "PATCH", "POST", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	if conf.Mode != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Static("/upload", "./upload")

	r.POST("/signup", h.Signup)
	r.POST("/login", h.Login)
	r.GET("/galleries", h.ListPublish)
	r.GET("/galleries/:id", h.GetOne)
	r.GET("/galleries/:id/images", h.ListGalleryImages)
	auth := r.Group("/")
	auth.Use(middleware.RequireUser(us))
	{
		auth.POST("/logout", h.Logout)
		user := auth.Group("/user")
		{
			user.POST("/galleries", h.Create)
			user.GET("/galleries", h.List)
			user.GET("/galleries/:id", h.GetOneAndCheck)
			user.DELETE("/galleries/:id", h.Delete)
			user.PATCH("/galleries/:id/names", h.UpdateName)
			user.PATCH("/galleries/:id/publishes", h.UpdatePublishing)
			user.POST("/galleries/:id/images", h.CreateImage)
			user.DELETE("/galleries/:id/images/:filename", h.DeleteImage)
			user.DELETE("/galleries/:id/images", h.DeleteImageInGallary)
			user.GET("/profile", h.GetProfile)
			user.POST("/profile", h.UpdateProfile)
		}
	}
	r.Run(":8080")
}

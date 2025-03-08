package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"warehouse-backend/database"
	"warehouse-backend/models"
	"warehouse-backend/routes"

	_ "warehouse-backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func runMigrations() error {
	database.Connect()
	db := database.GetDB()
	err := db.AutoMigrate(&models.User{}, &models.Product{})
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}
	fmt.Println("✅ Migrasi database berhasil!")
	return nil
}

// @title Simple Warehouse API
// @version 1.0
// @description API untuk mengelola gudang sederhana.
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Inisialisasi database
	database.Connect()
	if database.DB == nil {
		log.Fatal("Database connection is not initialized")
	}

	// Cek argumen CLI untuk migrasi
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		runMigrations()
		return
	}

	// Inisialisasi router
	r := gin.Default()

	// CORS configuration - place this BEFORE any routes
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"Accept",
			"Cache-Control",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Enable CORS for all routes using middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Cache-Control, X-Requested-With")

		// Handle OPTIONS method
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Inisialisasi Swagger API Documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup Routes
	routes.AuthRoutes(r)
	routes.ProductRoutes(r)

	// Server run on port 8080
	log.Println("Server running on port 8080")
	r.Run(":8080")
}

package main

import (
	"fmt"
	"log"
	"os"
	"warehouse-backend/database"
	"warehouse-backend/models"
	"warehouse-backend/routes"

	_ "warehouse-backend/docs"

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
	// Cek argumen CLI untuk migrasi
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		runMigrations()
		return
	}

	// Inisialisasi router
	r := gin.Default()

	// Inisialisasi Swagger API Documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup Routes
	routes.AuthRoutes(r)
	routes.ProductRoutes(r)

	// Server run on port 8080
	log.Println("Server running on port 8080")
	r.Run(":8080")
}

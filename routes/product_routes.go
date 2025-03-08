package routes

import (
	"warehouse-backend/controllers"
	"warehouse-backend/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.Engine) {
	productGroup := r.Group("/api/products")
	productGroup.Use(middleware.AuthMiddleware())
	{
		productGroup.POST("/", controllers.CreateProduct)
		productGroup.GET("/", controllers.GetProducts)

		productGroup.GET("/:id", controllers.GetProductByID)
		productGroup.PUT("/:id", controllers.UpdateProduct)
		productGroup.PUT("/:id/stock", controllers.UpdateStock)
		productGroup.DELETE("/:id", controllers.DeleteProduct)

		productGroup.GET("/barcode/:sku", controllers.GetBarcode)
		productGroup.GET("/export", controllers.ExportProductsCSV)
		productGroup.GET("/dashboard", controllers.GetStockDashboard)
		productGroup.POST("/bulk", controllers.BulkInsertProducts)
	}
}

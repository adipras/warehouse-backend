package controllers

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"warehouse-backend/database"
	"warehouse-backend/models"
	"warehouse-backend/utils"

	"github.com/gin-gonic/gin"
)

// BulkInsertProducts godoc
// @Summary Menambahkan banyak produk sekaligus
// @Description Menyimpan data produk dalam jumlah banyak menggunakan JSON array
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param products body []models.ProductSwagger true "Daftar Produk"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/bulk [post]
func BulkInsertProducts(c *gin.Context) {
	var products []models.Product

	// Parse JSON request body
	if err := c.ShouldBindJSON(&products); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi: Pastikan data tidak kosong
	if len(products) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data produk tidak boleh kosong"})
		return
	}

	// Insert ke database
	if err := database.DB.Create(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan produk"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Produk berhasil ditambahkan"})
}

type StockUpdateRequest struct {
	Change int `json:"change"`
}

// UpdateStock godoc
// @Summary Update stock of a product
// @Description Change the quantity of a product by a specified amount
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param body body StockUpdateRequest true "Stock change request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{id}/stock [put]
// @Security BearerAuth
func UpdateStock(c *gin.Context) {
	id := c.Param("id")
	var request StockUpdateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var product models.Product
	db := database.GetDB()
	if err := db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	product.Quantity += request.Change
	if product.Quantity < 0 {
		product.Quantity = 0
	}
	if product.Quantity == 0 {
		product.Status = "Out of Stock"
	} else if product.Quantity > 0 && product.Quantity < 10 {
		product.Status = "Low Stock"
	} else {
		product.Status = "Available"
	}

	if err := db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock updated", "quantity": product.Quantity})
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the input payload
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body models.ProductSwagger true "Product JSON"
// @Success 201 {object} models.CreateProductResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	product.Status = "Available"
	if product.Quantity == 0 {
		product.Status = "Out of Stock"
	} else if product.Quantity > 0 && product.Quantity < 10 {
		product.Status = "Low Stock"
	}

	// Generate SKU otomatis
	product.SKU = utils.GenerateSKU()

	// Generate barcode
	barcodePath, err := utils.GenerateBarcode(product.SKU)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat barcode"})
		return
	}

	product.BarcodePath = barcodePath

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, models.CreateProductResponse{Message: "Product created successfully"})
}

// GetProducts godoc
// @Summary Get all products
// @Description Get all products
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.ProductSwagger
// @Router /products [get]
func GetProducts(c *gin.Context) {
	var products []models.Product
	database.DB.Find(&products)
	c.JSON(http.StatusOK, products)
}

// GetProductByID godoc
// @Summary Get a product by ID
// @Description Get a product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} models.ProductSwagger
// @Failure 404 {object} models.ErrorResponse
// @Router /products/{id} [get]
func GetProductByID(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetBarcode godoc
// @Summary Ambil barcode produk
// @Description Mengembalikan gambar barcode berdasarkan SKU
// @Tags Products
// @Produce png
// @Security BearerAuth
// @Param sku path string true "Product SKU"
// @Success 200
// @Failure 404 {object} map[string]string
// @Router /products/barcode/{sku} [get]
func GetBarcode(c *gin.Context) {
	sku := c.Param("sku")

	var product models.Product
	if err := database.DB.Where("sku = ?", sku).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	c.File(product.BarcodePath)
}

// ExportProductsCSV godoc
// @Summary Ekspor daftar produk ke CSV
// @Description Mengunduh daftar produk dalam format CSV
// @Tags Products
// @Accept json
// @Produce text/csv
// @Security BearerAuth
// @Success 200 {string} csv "File CSV"
// @Failure 500 {object} map[string]string
// @Router /products/export [get]
func ExportProductsCSV(c *gin.Context) {
	var products []models.Product
	database.DB.Find(&products)

	// Set header untuk file CSV
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=products.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Menulis header CSV
	writer.Write([]string{"ID", "Name", "SKU", "Quantity", "Location", "Status", "BarcodePath"})

	// Menulis data produk ke CSV
	for _, product := range products {
		writer.Write([]string{
			strconv.Itoa(int(product.ID)),
			product.Name,
			product.SKU,
			strconv.Itoa(product.Quantity),
			product.Location,
			product.Status,
			product.BarcodePath,
		})
	}
}

// GetStockDashboard godoc
// @Summary Mendapatkan ringkasan stok gudang
// @Description Mengambil total stok, jumlah produk, dan daftar produk dengan stok rendah
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/dashboard [get]
func GetStockDashboard(c *gin.Context) {
	var totalProducts int64
	var outOfStock int64
	var lowStock int64

	database.DB.Model(&models.Product{}).Count(&totalProducts)
	database.DB.Model(&models.Product{}).Where("quantity = 0").Count(&outOfStock)
	database.DB.Model(&models.Product{}).Where("quantity < ?", 10).Count(&lowStock)

	var latestProducts []models.Product
	database.DB.Order("created_at DESC").Limit(5).Find(&latestProducts)

	c.JSON(http.StatusOK, gin.H{
		"summary": gin.H{
			"totalProducts": totalProducts,
			"outOfStock":    outOfStock,
			"lowStock":      lowStock,
		},
		"latestProducts": latestProducts,
	})
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update a product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param product body models.ProductSwagger true "Product JSON"
// @Success 200 {object} models.CreateProductResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product not found"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	database.DB.Save(&product)
	c.JSON(http.StatusOK, models.CreateProductResponse{Message: "Product updated successfully"})
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} models.DeleteProductResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product not found"})
		return
	}

	database.DB.Delete(&product)
	c.JSON(http.StatusOK, models.DeleteProductResponse{Message: "Product deleted successfully"})
}

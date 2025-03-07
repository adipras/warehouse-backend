package models

import (
	"gorm.io/gorm"
)

// Product represents a product in the warehouse
type Product struct {
	gorm.Model         // ID, CreatedAt, UpdatedAt, DeletedAt
	Name        string `gorm:"type:varchar(255);not null" json:"name" example:"Produk A"`
	SKU         string `gorm:"type:varchar(100);uniqueIndex;not null" json:"sku" example:"SKU123"`
	Quantity    int    `gorm:"not null" json:"quantity" example:"100"`
	Location    string `gorm:"type:varchar(255)" json:"location" example:"Rak 1"`
	Status      string `gorm:"type:varchar(50);not null" json:"status" example:"available"`
	BarcodePath string `json:"barcode_path"`
}

// ProductSwagger represents a product in the warehouse for Swagger documentation
// @Description Product represents a product in the warehouse
type ProductSwagger struct {
	Name     string `json:"name" example:"Produk A"`
	Quantity int    `json:"quantity" example:"100"`
	Location string `json:"location" example:"Rak 1"`
	Status   string `json:"status" example:"available"`
}

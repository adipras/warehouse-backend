package utils

import (
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/boombuler/barcode/code128"
	"golang.org/x/image/draw"
)

// GenerateBarcode membuat barcode dari SKU dan menyimpannya sebagai file PNG
func GenerateBarcode(sku string) (string, error) {
	barcodePath := "storage/barcodes/"
	os.MkdirAll(barcodePath, os.ModePerm)

	// Buat barcode dari SKU
	barcodeData, err := code128.Encode(sku)
	if err != nil {
		return "", err
	}

	// Get the original bounds
	bounds := barcodeData.Bounds()

	// Create a new RGBA image with desired dimensions
	newImage := image.NewRGBA(image.Rect(0, 0, 300, 150))

	// Scale the image using draw.ApproxBiLinear
	draw.ApproxBiLinear.Scale(newImage, newImage.Bounds(), barcodeData, bounds, draw.Over, nil)

	// Path file barcode
	filePath := filepath.Join(barcodePath, sku+".png")

	// Simpan barcode sebagai PNG
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = png.Encode(file, newImage)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// DeleteBarcode removes the barcode image file for a given SKU
func DeleteBarcode(sku string) error {
	barcodePath := filepath.Join("storage/barcodes", sku+".png")

	// Check if file exists before attempting to delete
	if _, err := os.Stat(barcodePath); os.IsNotExist(err) {
		return nil // Return nil if file doesn't exist
	}

	// Delete the barcode file
	err := os.Remove(barcodePath)
	if err != nil {
		return err
	}

	return nil
}

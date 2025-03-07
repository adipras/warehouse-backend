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

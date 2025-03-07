package utils

import (
	"fmt"
	"time"
)

// GenerateSKU membuat SKU berdasarkan timestamp
func GenerateSKU() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond) // Timestamp dalam milidetik
	return fmt.Sprintf("SKU-%d", timestamp)
}

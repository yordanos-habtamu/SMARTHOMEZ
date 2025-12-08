package utils

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

// GenerateQRCode creates a QR code image and uploads it to Cloudinary
// Returns the Cloudinary URL of the QR code image
func GenerateQRCode(content string) (string, error) {
	// Generate QR code as PNG bytes
	qrCode, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %v", err)
	}

	// Upload to Cloudinary
	url, err := UploadQRCodeToCloudinary(qrCode)
	if err != nil {
		return "", fmt.Errorf("failed to upload QR code: %v", err)
	}

	return url, nil
}

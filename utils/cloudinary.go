package utils

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/yordanos-habtamu/realstate/config"
)

var cld *cloudinary.Cloudinary

// InitCloudinary initializes the Cloudinary client
func InitCloudinary() error {
	var err error
	cld, err = cloudinary.NewFromParams(
		config.Envs.CLOUDINARY_CLOUD_NAME,
		config.Envs.CLOUDINARY_API_KEY,
		config.Envs.CLOUDINARY_API_SECRET,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize Cloudinary: %v", err)
	}
	return nil
}

// UploadHouseImage uploads a house image to Cloudinary
func UploadHouseImage(file multipart.File, filename string) (string, error) {
	if cld == nil {
		if err := InitCloudinary(); err != nil {
			return "", err
		}
	}

	ctx := context.Background()
	uploadParams := uploader.UploadParams{
		Folder:         "realstate/images",
		PublicID:       filename,
		ResourceType:   "image",
		Transformation: "c_limit,w_1200,h_900,q_auto:good",
	}

	result, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %v", err)
	}

	return result.SecureURL, nil
}

// UploadQRCodeToCloudinary uploads a QR code image to Cloudinary
func UploadQRCodeToCloudinary(qrCodeBytes []byte) (string, error) {
	if cld == nil {
		if err := InitCloudinary(); err != nil {
			return "", err
		}
	}

	ctx := context.Background()

	// Create a reader from bytes
	reader := bytes.NewReader(qrCodeBytes)

	uploadParams := uploader.UploadParams{
		Folder:       "qrcodes",
		ResourceType: "image",
	}

	result, err := cld.Upload.Upload(ctx, reader, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload QR code: %v", err)
	}

	return result.SecureURL, nil
}

// DeleteCloudinaryImage deletes an image from Cloudinary using its public ID
func DeleteCloudinaryImage(publicID string) error {
	if cld == nil {
		if err := InitCloudinary(); err != nil {
			return err
		}
	}

	ctx := context.Background()
	_, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})

	if err != nil {
		return fmt.Errorf("failed to delete image: %v", err)
	}

	return nil
}

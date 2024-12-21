package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/services"
)

// UploadImageHandler handles image uploads to Cloudinary
func UploadImageHandler(c *fiber.Ctx) error {
	// Validate user token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required. Please log in.",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token. Please log in again.",
		})
	}

	// Initialize Cloudinary client
	cld := services.CloudinaryClient
	ctx := context.Background()

	// Parse the uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error retrieving file from request",
		})
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to open file",
		})
	}
	defer file.Close()

	// Upload to Cloudinary
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: fileHeader.Filename,
		Folder:   fmt.Sprintf("users/%s/images", uid),
	})
	if err != nil {
		log.Println("Error uploading image:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload image. Please try again.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Image uploaded successfully",
		"imageUrl": uploadResult.SecureURL,
	})
}

// UploadPDFHandler handles PDF uploads to Cloudinary
func UploadPDFHandler(c *fiber.Ctx) error {
	// Validate user token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required. Please log in.",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token. Please log in again.",
		})
	}

	// Initialize Cloudinary client
	cld := services.CloudinaryClient
	ctx := context.Background()

	// Parse the uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error retrieving file from request",
		})
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to open file",
		})
	}
	defer file.Close()

	// Upload the file to Cloudinary
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:     fileHeader.Filename,
		Folder:       fmt.Sprintf("users/%s/pdfs", uid),
		ResourceType: "raw", // 'raw' is used for non-image and non-video files like PDFs
	})
	if err != nil {
		log.Println("Error uploading PDF:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload PDF. Please try again.",
		})
	}

	// Return a success response with the PDF URL
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "PDF uploaded successfully",
		"pdf_url": uploadResult.SecureURL,
	})
}

// UploadVideoHandler handles video uploads to Cloudinary
func UploadVideoHandler(c *fiber.Ctx) error {
	// Validate user token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required. Please log in.",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token. Please log in again.",
		})
	}

	// Initialize Cloudinary client
	cld := services.CloudinaryClient
	ctx := context.Background()

	// Parse the uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error retrieving file from request",
		})
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to open file",
		})
	}
	defer file.Close()

	// Upload the file to Cloudinary
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:     fileHeader.Filename,
		Folder:       fmt.Sprintf("users/%s/videos", uid),
		ResourceType: "video", // Specify resource type for videos
	})
	if err != nil {
		log.Println("Error uploading video:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload video. Please try again.",
		})
	}

	// Return a success response with the video URL
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Video uploaded successfully",
		"video_url": uploadResult.SecureURL,
	})
}

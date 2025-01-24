package services

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/rogerjeasy/go-letusconnect/config"
)

var CloudinaryClient *cloudinary.Cloudinary

func InitCloudinary() *cloudinary.Cloudinary {
	cloudinaryURL := config.CloudinaryURL
	if cloudinaryURL == "" {
		log.Fatal("Cloudinary URL not found in environment variables")
	}
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}
	CloudinaryClient = cld
	log.Println("Cloudinary client initialized successfully")
	return cld
}

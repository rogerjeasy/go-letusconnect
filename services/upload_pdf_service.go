package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type UploadPDFService struct {
	firestoreClient FirestoreClient
	cloudinary      *cloudinary.Cloudinary
}

type UploadResponse struct {
	URL     string `json:"url"`
	Message string `json:"message"`
}

func NewUploadPDFService(firestoreClient FirestoreClient, cloudinaryURL string) (*UploadPDFService, error) {
	// Initialize Cloudinary
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, fmt.Errorf("error initializing Cloudinary: %v", err)
	}

	return &UploadPDFService{
		firestoreClient: firestoreClient,
		cloudinary:      cld,
	}, nil
}

func (s *UploadPDFService) UploadPDF(file *multipart.FileHeader) (*UploadResponse, error) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "upload-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Save the uploaded file to temp location
	sourceFile, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening uploaded file: %v", err)
	}
	defer sourceFile.Close()

	// Copy the uploaded file to the temporary file
	buffer := make([]byte, 1024)
	for {
		n, err := sourceFile.Read(buffer)
		if err != nil {
			break
		}
		tempFile.Write(buffer[:n])
	}

	// Upload to Cloudinary
	ctx := context.Background()
	uploadResult, err := s.cloudinary.Upload.Upload(
		ctx,
		tempFile.Name(),
		uploader.UploadParams{
			ResourceType: "raw",
			Folder:       "website-context",
			PublicID:     "project-description-" + file.Filename,
			Type:         "upload",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error uploading to Cloudinary: %v", err)
	}

	// Store upload record in Firestore if needed
	if err := s.storeUploadRecord(ctx, uploadResult.SecureURL, file.Filename); err != nil {
		// Log error but don't fail the upload
		fmt.Printf("Error storing upload record: %v\n", err)
	}

	return &UploadResponse{
		URL:     uploadResult.SecureURL,
		Message: "PDF uploaded successfully to Cloudinary",
	}, nil
}

func (s *UploadPDFService) storeUploadRecord(ctx context.Context, url string, filename string) error {
	_, _, err := s.firestoreClient.Collection("website-context").Add(ctx, map[string]interface{}{
		"url":        url,
		"filename":   filename,
		"uploadedAt": firestore.ServerTimestamp,
		"type":       "website-context",
	})
	return err
}

package models

import (
	"time"
)

type Attachment struct {
	FileName   string    `json:"file_name"`
	URL        string    `json:"url"`
	UploadedAt time.Time `json:"uploaded_at"`
}

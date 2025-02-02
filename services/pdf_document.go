package services

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ledongthuc/pdf"
)

type PDFService struct {
	firestoreClient FirestoreClient
	context         string
	mutex           sync.RWMutex
	lastUpdated     time.Time
	pdfURL          string
	stopRefresh     chan struct{}
}

func NewPDFService(firestoreClient FirestoreClient, pdfURL string) *PDFService {
	if pdfURL == "" {
		log.Fatal("PDF URL cannot be empty")
	}

	service := &PDFService{
		firestoreClient: firestoreClient,
		pdfURL:          pdfURL,
		stopRefresh:     make(chan struct{}),
	}

	// Initial load of the context
	if err := service.RefreshContext(); err != nil {
		log.Printf("Initial PDF context load failed: %v", err)
	}

	// Start periodic refresh
	go service.startPeriodicRefresh()

	return service
}

func (s *PDFService) startPeriodicRefresh() {
	ticker := time.NewTicker(24 * time.Hour) // Refresh every 24 hours
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.RefreshContext(); err != nil {
				log.Printf("Failed to refresh PDF context: %v", err)
			}
		case <-s.stopRefresh:
			return
		}
	}
}

func (s *PDFService) Stop() {
	close(s.stopRefresh)
}

func (s *PDFService) GetContext() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.context
}

func (s *PDFService) RefreshContext() error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Download PDF from URL
	resp, err := client.Get(s.pdfURL)
	if err != nil {
		return fmt.Errorf("error downloading PDF: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Extract text from PDF
	content, err := s.extractPDFText(resp.Body)
	if err != nil {
		return fmt.Errorf("error extracting PDF text: %v", err)
	}

	if content == "" {
		return fmt.Errorf("extracted PDF content is empty")
	}

	// Update the context
	s.mutex.Lock()
	s.context = content
	s.lastUpdated = time.Now()
	s.mutex.Unlock()

	return nil
}

func (s *PDFService) GetLastUpdated() time.Time {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastUpdated
}

func (s *PDFService) extractPDFText(reader io.Reader) (string, error) {
	tempFile, err := os.CreateTemp("", "pdf-*.pdf")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	// Copy the reader content to the temp file
	if _, err := io.Copy(tempFile, reader); err != nil {
		return "", fmt.Errorf("failed to copy content to temp file: %v", err)
	}

	// Close the file for writing and reopen for reading
	if err := tempFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %v", err)
	}

	// Open the PDF file for reading
	f, r, err := pdf.Open(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %v", err)
	}
	defer f.Close()

	var buffer bytes.Buffer
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		text, err := p.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("failed to extract text from page %d: %v", pageIndex, err)
		}
		buffer.WriteString(text)
		buffer.WriteString("\n")
	}

	// Clean and normalize the text
	text := buffer.String()
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	if text == "" {
		return "", fmt.Errorf("no text content extracted from PDF")
	}

	return text, nil
}

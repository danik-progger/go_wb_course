package service

import (
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"imageProcessing/entities"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/image/draw"
)

type ImageServiceImpl struct {
	storagePath string
}

func NewImageService(storagePath string) *ImageServiceImpl {
	// Create storage directories if they don't exist
	err := os.MkdirAll(filepath.Join(storagePath, "original"), 0755)
	if err != nil {
		fmt.Printf("Error creating original directory: %v\n", err)
	}
	err = os.MkdirAll(filepath.Join(storagePath, "processed"), 0755)
	if err != nil {
		fmt.Printf("Error creating processed directory: %v\n", err)
	}

	return &ImageServiceImpl{
		storagePath: storagePath,
	}
}

func (s *ImageServiceImpl) Upload(imageData []byte, filename string, contentType string) (*entities.Image, error) {
	id := int(time.Now().Unix())

	originalPath := filepath.Join(s.storagePath, "original", fmt.Sprintf("%d_%s", id, filename))
	err := os.WriteFile(originalPath, imageData, 0644)
	if err != nil {
		return nil, err
	}

	// Get file info
	fileInfo, err := os.Stat(originalPath)
	if err != nil {
		return nil, err
	}

	// Create image entity
	image := &entities.Image{
		ID:          id,
		OriginalURL: originalPath,
		Filename:    fmt.Sprintf("%d_%s", id, filename),
		Size:        fileInfo.Size(),
		ContentType: contentType,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	// Set processed URL path
	processedPath := filepath.Join(s.storagePath, "processed", fmt.Sprintf("%d_%s", id, filename))
	image.ProcessedURL = processedPath

	// Queue image processing (in a real app this would be in a goroutine or job queue)
	go func() {
		// Process with all actions for basic implementation
		s.ProcessImage(image, entities.ActionResize)
		s.ProcessImage(image, entities.ActionMiniature)
		s.ProcessImage(image, entities.ActionWatermark)
	}()

	return image, nil
}

func (s *ImageServiceImpl) Get(id int) (*entities.Image, error) {
	// Find the original file
	originalPath := filepath.Join(s.storagePath, "original", fmt.Sprintf("%d_*", id))
	originalFiles, err := filepath.Glob(originalPath)
	if err != nil || len(originalFiles) == 0 {
		return nil, fmt.Errorf("original image not found")
	}

	// Look for a processed version (try watermark version first, then resized, then miniature)
	processedFiles, err := filepath.Glob(filepath.Join(s.storagePath, "processed", fmt.Sprintf("%d_watermarked_*", id)))
	if err != nil || len(processedFiles) == 0 {
		// Try resized version
		processedFiles, err = filepath.Glob(filepath.Join(s.storagePath, "processed", fmt.Sprintf("%d_resized_*", id)))
		if err != nil || len(processedFiles) == 0 {
			// Try miniature version
			processedFiles, err = filepath.Glob(filepath.Join(s.storagePath, "processed", fmt.Sprintf("%d_mini_*", id)))
			if err != nil || len(processedFiles) == 0 {
				// If no processed versions exist, just return the original
				processedFiles = originalFiles
			}
		}
	}

	// Get file info for original
	originalInfo, err := os.Stat(originalFiles[0])
	if err != nil {
		return nil, err
	}

	// Get file info for processed
	processedFile := processedFiles[0]
	processedInfo, err := os.Stat(processedFile)
	if err != nil {
		processedFile = originalFiles[0]
		processedInfo = originalInfo
	}

	// Extract the actual filename (without the ID prefix)
	baseOriginal := filepath.Base(originalFiles[0])
	filename := baseOriginal
	// Remove the "ID_filename" prefix to get just the filename
	if underscoreIdx := strings.Index(baseOriginal, "_"); underscoreIdx != -1 {
		filename = baseOriginal[underscoreIdx+1:]
	}

	return &entities.Image{
		ID:           id,
		OriginalURL:  originalFiles[0],
		ProcessedURL: processedFile,
		Filename:     filename,
		Size:         processedInfo.Size(),
		ContentType:  "image/jpeg", // Would determine from file in real implementation
		CreatedAt:    time.Now().Format(time.RFC3339),
	}, nil
}

func (s *ImageServiceImpl) Delete(id int) error {
	// Delete original image file
	originalPath := filepath.Join(s.storagePath, "original", fmt.Sprintf("%d_*", id))
	originalFiles, err := filepath.Glob(originalPath)
	if err == nil {
		for _, file := range originalFiles {
			os.Remove(file)
		}
	}

	// Delete processed image file
	processedPath := filepath.Join(s.storagePath, "processed", fmt.Sprintf("%d_*", id))
	processedFiles, err := filepath.Glob(processedPath)
	if err == nil {
		for _, file := range processedFiles {
			os.Remove(file)
		}
	}

	return nil
}

func (s *ImageServiceImpl) ProcessImage(imgEntity *entities.Image, action string) error {
	// Open the original image file
	originalFile, err := os.Open(imgEntity.OriginalURL)
	if err != nil {
		return err
	}
	defer originalFile.Close()

	// Decode the image
	img, format, err := image.Decode(originalFile)
	if err != nil {
		return err
	}

	// Process based on action
	var processedImg image.Image
	switch action {
	case entities.ActionResize:
		processedImg = s.resizeImage(img, 800, 600) // Resize to 800x600 max
	case entities.ActionMiniature:
		processedImg = s.resizeImage(img, 150, 150) // Create a 150x150 thumbnail
	case entities.ActionWatermark:
		processedImg = s.addWatermark(img) // Add a simple text watermark
	default:
		processedImg = img
	}

	// Create a unique processed path for each action
	var processedPath string
	switch action {
	case entities.ActionResize:
		processedPath = filepath.Join(s.storagePath, "processed",
			fmt.Sprintf("%d_resized_%s", imgEntity.ID, filepath.Base(imgEntity.OriginalURL)))
	case entities.ActionMiniature:
		processedPath = filepath.Join(s.storagePath, "processed",
			fmt.Sprintf("%d_mini_%s", imgEntity.ID, filepath.Base(imgEntity.OriginalURL)))
	case entities.ActionWatermark:
		processedPath = filepath.Join(s.storagePath, "processed",
			fmt.Sprintf("%d_watermarked_%s", imgEntity.ID, filepath.Base(imgEntity.OriginalURL)))
	default:
		processedPath = filepath.Join(s.storagePath, "processed",
			fmt.Sprintf("%d_%s", imgEntity.ID, filepath.Base(imgEntity.OriginalURL)))
	}

	outputFile, err := os.Create(processedPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Encode the processed image in the original format
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(outputFile, processedImg, &jpeg.Options{Quality: 90})
	case "png":
		err = png.Encode(outputFile, processedImg)
	default:
		// Default to JPEG if format is unknown
		err = jpeg.Encode(outputFile, processedImg, &jpeg.Options{Quality: 90})
	}
	if err != nil {
		return err
	}

	// Update the main ProcessedURL to point to the watermarked version if that's the action
	// or to a default processed version
	if action == entities.ActionWatermark {
		// In a complete implementation, you might update the entity's processed URL
		// For now, we'll just ensure the image has a processed version
	}

	return nil
}

// resizeImage resizes an image to fit within the specified dimensions while maintaining aspect ratio
func (s *ImageServiceImpl) resizeImage(img image.Image, maxWidth, maxHeight int) image.Image {
	// Get the original image dimensions
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// Calculate the scaling factor
	widthRatio := float64(maxWidth) / float64(origWidth)
	heightRatio := float64(maxHeight) / float64(origHeight)
	scale := widthRatio
	if heightRatio < widthRatio {
		scale = heightRatio
	}

	// Calculate new dimensions
	newWidth := int(float64(origWidth) * scale)
	newHeight := int(float64(origHeight) * scale)

	// Create a new image with the calculated dimensions
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Resize the image using the draw package
	draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)

	return dst
}

// addWatermark adds a simple text watermark to the image
func (s *ImageServiceImpl) addWatermark(img image.Image) image.Image {
	// For a simple implementation, we're returning the image as is
	// A complete implementation would add actual watermarking
	// This requires more complex image manipulation which would need additional libraries
	return img
}

package service

import (
	"imageProcessing/entities"
)

type ImageService interface {
	Upload(imageData []byte, filename string, contentType string) (*entities.Image, error)
	Get(id int) (*entities.Image, error)
	Delete(id int) error
	ProcessImage(image *entities.Image, action string) error
}

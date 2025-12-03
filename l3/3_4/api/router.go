package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"imageProcessing/service"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router       *chi.Mux
	imageService service.ImageService
}

func NewServer(imageService service.ImageService) *Server {
	s := &Server{
		Router:       chi.NewRouter(),
		imageService: imageService,
	}

	// Setup routes
	s.Router.Get("/", s.serveFrontend)
	s.Router.Post("/upload", s.uploadImageHandler)
	s.Router.Get("/image/{id}", s.getImageHandler)
	s.Router.Delete("/image/{id}", s.deleteImageHandler)

	return s
}

func (s *Server) uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file content
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}

	// Upload image
	image, err := s.imageService.Upload(fileData, fileHeader.Filename, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, "Error uploading image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(image)
}

func (s *Server) getImageHandler(w http.ResponseWriter, r *http.Request) {
	// Get image ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	// Get image
	image, err := s.imageService.Get(id)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	// Serve the processed image file
	http.ServeFile(w, r, image.ProcessedURL)
}

func (s *Server) deleteImageHandler(w http.ResponseWriter, r *http.Request) {
	// Get image ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	// Delete image
	err = s.imageService.Delete(id)
	if err != nil {
		http.Error(w, "Error deleting image", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Image deleted successfully")
}

func (s *Server) serveFrontend(w http.ResponseWriter, r *http.Request) {
	// Try multiple possible paths for the frontend file
	paths := []string{
		"./front/index.html",
		"front/index.html",
		"../front/index.html",
		"../../front/index.html",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			http.ServeFile(w, r, path)
			return
		}
	}

	// If no file found, return error
	http.Error(w, "Frontend file not found", http.StatusNotFound)
}

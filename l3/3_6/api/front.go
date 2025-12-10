package api

import (
	"net/http"
	"os"
	"salesTracker/db"

	"github.com/go-chi/chi/v5"
)

type FrontendServer struct {
	Router *chi.Mux
}

func NewFrontendServer(salesDB *db.DB) *FrontendServer {
	s := &FrontendServer{
		Router: chi.NewRouter(),
	}

	// Setup routes
	s.Router.Get("/", s.serveFrontend)

	return s
}

func (s *FrontendServer) serveFrontend(w http.ResponseWriter, r *http.Request) {
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

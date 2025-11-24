package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"commenttree/entities"

	"github.com/go-chi/chi/v5"
)

type CommentService interface {
	AddComment(parentID *int, text string) *entities.Comment
	DeleteComment(id int) bool
	SearchComments(query string) []*entities.Comment
	GetCommentsByParent(parentID *int) []*entities.Comment
	GetCommentWithSubtree(commentID int) *entities.Comment
}

type Server struct {
	Router         *chi.Mux
	commentService CommentService
}

func NewServer(commentService CommentService) *Server {
	s := &Server{
		Router:         chi.NewRouter(),
		commentService: commentService,
	}

	// Setup routes
	s.Router.Post("/comments", s.createCommentHandler)
	s.Router.Get("/comments", s.getCommentsHandler)
	s.Router.Delete("/comments/{id}", s.deleteCommentHandler)
	s.Router.Get("/comments/search", s.searchCommentsHandler)

	return s
}

// Handler for creating a new comment
func (s *Server) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var req struct {
		ParentID *int   `json:"parent_id"`
		Text     string `json:"text"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	if req.Text == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	comment := s.commentService.AddComment(req.ParentID, req.Text)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

// Handler for getting comments
func (s *Server) getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// Add a special header to identify if this handler is being used
	w.Header().Set("X-Handler", "getCommentsHandler-updated")

	parentIDStr := r.URL.Query().Get("parent")
	if parentIDStr != "" {
		// If parent is specified, return the comment with that ID and its full subtree
		id, err := strconv.Atoi(parentIDStr)
		if err != nil {
			http.Error(w, "Invalid parent ID", http.StatusBadRequest)
			return
		}

		comment := s.commentService.GetCommentWithSubtree(id)
		if comment == nil {
			http.Error(w, "Comment not found - FROM UPDATED HANDLER", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comment)
	} else {
		// If no parent is specified, return all root comments with their subtrees
		comments := s.commentService.GetCommentsByParent(nil)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	}
}

// Handler for deleting a comment
func (s *Server) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	if !s.commentService.DeleteComment(id) {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Comment deleted")
}

// Handler for searching comments
func (s *Server) searchCommentsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	comments := s.commentService.SearchComments(query)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

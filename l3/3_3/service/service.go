package service

import (
	"strings"
	"sync"
	"time"

	"commenttree/entities"
)

// CommentService manages entities.Comments storage and operations
type CommentService struct {
	Comments map[int]*entities.Comment
	mutex    sync.RWMutex
	nextID   int
}

// NewCommentService creates a new CommentService instance
func NewCommentService() *CommentService {
	return &CommentService{
		Comments: make(map[int]*entities.Comment),
		nextID:   1,
	}
}

// AddComment creates a new Comment
func (cs *CommentService) AddComment(parentID *int, text string) *entities.Comment {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	comment := &entities.Comment{
		ID:        cs.nextID,
		ParentID:  parentID,
		Text:      text,
		CreatedAt: time.Now(),
		Children:  []*entities.Comment{},
	}
	cs.Comments[cs.nextID] = comment
	cs.nextID++
	return comment
}

// GetCommentsByParent returns Comments with a specific parent (or root entities.Comments if parentID is nil)
// When parentID is nil, returns all root comments each with their full subtree
// When parentID is not nil, returns direct children of that parent (not the parent itself with subtree)
func (cs *CommentService) GetCommentsByParent(parentID *int) []*entities.Comment {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	var result []*entities.Comment
	for _, comment := range cs.Comments {
		if parentID == nil && comment.ParentID == nil {
			// Root entities.Comment - return with full subtree
			result = append(result, cs.buildTree(comment))
		} else if parentID != nil && comment.ParentID != nil && *comment.ParentID == *parentID {
			// Direct children of specified parent
			result = append(result, cs.buildTree(comment))
		}
	}
	return result
}

// GetCommentWithSubtree returns a specific comment with its full subtree
func (cs *CommentService) GetCommentWithSubtree(commentID int) *entities.Comment {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	comment, exists := cs.Comments[commentID]
	if !exists {
		return nil
	}

	// Return the comment with its full subtree
	return cs.buildTree(comment)
}

// buildTree recursively builds the comment tree
func (cs *CommentService) buildTree(comment *entities.Comment) *entities.Comment {
	comment.Children = cs.getChildren(comment.ID)
	for _, child := range comment.Children {
		cs.buildTree(child)
	}
	return comment
}

// getChildren returns all direct children of a entities.Comment
func (cs *CommentService) getChildren(parentID int) []*entities.Comment {
	var children []*entities.Comment
	for _, comment := range cs.Comments {
		if comment.ParentID != nil && *comment.ParentID == parentID {
			children = append(children, comment)
		}
	}
	return children
}

// DeleteComment deletes a comment and all its children recursively
func (cs *CommentService) DeleteComment(id int) bool {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	if _, exists := cs.Comments[id]; !exists {
		return false
	}

	// Delete all children first
	cs.deleteChildren(id)

	// Delete the entities.Comment itself
	delete(cs.Comments, id)
	return true
}

// deleteChildren recursively deletes all children of a entities.Comment
func (cs *CommentService) deleteChildren(parentID int) {
	for id, comment := range cs.Comments {
		if comment.ParentID != nil && *comment.ParentID == parentID {
			// Delete children of this child first
			cs.deleteChildren(id)
			// Then delete the child
			delete(cs.Comments, id)
		}
	}
}

// SearchComments searches for comments containing the query text
func (cs *CommentService) SearchComments(query string) []*entities.Comment {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	var results []*entities.Comment
	query = strings.ToLower(query)
	for _, comment := range cs.Comments {
		if strings.Contains(strings.ToLower(comment.Text), query) {
			// Build tree for this entities.Comment and its context
			results = append(results, cs.buildTree(comment))
		}
	}
	return results
}

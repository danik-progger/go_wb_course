package entities

import "time"

type Comment struct {
	ID        int        `json:"id"`
	ParentID  *int       `json:"parent_id,omitempty"`
	Text      string     `json:"text"`
	CreatedAt time.Time  `json:"created_at"`
	Children  []*Comment `json:"children"`
}

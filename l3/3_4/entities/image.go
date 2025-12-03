package entities

type Image struct {
	ID           int    `json:"id"`
	OriginalURL  string `json:"original_url"`
	ProcessedURL string `json:"processed_url"`
	Filename     string `json:"filename"`
	Size         int64  `json:"size"`
	ContentType  string `json:"content_type"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	CreatedAt    string `json:"created_at"`
}

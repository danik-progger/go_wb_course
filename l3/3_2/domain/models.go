package domain

import "time"

type ShortenedURL struct {
	ID          string
	OriginalURL string
	CreatedAt   time.Time
}

type Visit struct {
	ShortURLID string
	Timestamp  time.Time
	UserAgent  string
}

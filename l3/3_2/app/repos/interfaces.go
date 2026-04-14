package repos

import (
	"urlshortener/domain"
)

type URLRepo interface {
	Create(url domain.ShortenedURL) error
	Get(id string) (domain.ShortenedURL, bool)
	Exists(id string) bool
}

type VisitRepo interface {
	Record(visit domain.Visit)
	GetByURL(id string) []domain.Visit
}

type CacheRepo interface {
	Get(key string) (string, bool)
	Set(key, value string)
}

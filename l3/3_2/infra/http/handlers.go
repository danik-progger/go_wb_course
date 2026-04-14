package http

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"urlshortener/app/service"
	"urlshortener/infra/templates"
)

type Handler struct {
	service *service.URLService
	baseURL string
}

func NewHandler(svc *service.URLService, baseURL string) *Handler {
	return &Handler{
		service: svc,
		baseURL: baseURL,
	}
}

func (h *Handler) HandleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/shorten", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, templates.FormTemplate)
}

func (h *Handler) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	customKey := r.FormValue("custom")

	shortKey, err := h.service.CreateShortURL(originalURL, customKey)
	if err != nil {
		if strings.Contains(err.Error(), "already taken") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	shortenedURL := fmt.Sprintf("%s/s/%s", h.baseURL, shortKey)
	analyticsURL := fmt.Sprintf("%s/analytics/%s", h.baseURL, shortKey)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, templates.ShortenResultTemplate,
		template.HTMLEscapeString(originalURL),
		template.HTMLEscapeString(shortenedURL),
		template.HTMLEscapeString(shortenedURL),
		template.HTMLEscapeString(analyticsURL),
		template.HTMLEscapeString(analyticsURL),
	)
}

func (h *Handler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := strings.TrimPrefix(r.URL.Path, "/s/")
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	originalURL, err := h.service.Resolve(shortKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.service.RecordVisit(shortKey, r.UserAgent())

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func (h *Handler) HandleAnalytics(w http.ResponseWriter, r *http.Request) {
	shortKey := strings.TrimPrefix(r.URL.Path, "/analytics/")
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	urlData, found := h.service.GetURL(shortKey)
	if !found {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		return
	}

	visits := h.service.GetVisits(shortKey)

	dailyCounts := make(map[string]int)
	monthlyCounts := make(map[string]int)
	uaCounts := make(map[string]int)

	for _, visit := range visits {
		day := visit.Timestamp.Format("2006-01-02")
		month := visit.Timestamp.Format("2006-01")
		dailyCounts[day]++
		monthlyCounts[month]++
		uaCounts[visit.UserAgent]++
	}

	data := templates.AnalyticsData{
		URL:             urlData,
		TotalVisits:     len(visits),
		Visits:          visits,
		DailyCounts:     dailyCounts,
		MonthlyCounts:   monthlyCounts,
		UserAgentCounts: uaCounts,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := templates.ExecuteAnalytics(w, data); err != nil {
		http.Error(w, "Failed to render analytics", http.StatusInternalServerError)
	}
}

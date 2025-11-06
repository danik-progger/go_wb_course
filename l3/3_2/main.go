package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Data structures
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

// In-memory stores
var (
	urlStore   = make(map[string]ShortenedURL)
	visitStore = make(map[string][]Visit)
	urlMutex   = &sync.RWMutex{}
	visitMutex = &sync.RWMutex{}
	// In-memory cache to simulate Redis
	urlCache   = make(map[string]string)
	cacheMutex = &sync.RWMutex{}
)

const (
	BaseURL = "http://localhost:3030"
)

func main() {
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/s/", handleRedirect)
	http.HandleFunc("/analytics/", handleAnalytics)

	fmt.Println("URL Shortener is running on :3030")
	http.ListenAndServe(":3030", nil)
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/shorten", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>URL Shortener</title>
			<style>
				body { font-family: sans-serif; max-width: 600px; margin: 40px auto; padding: 20px; }
				input[type=url], input[type=text] { width: 100%; padding: 8px; margin-bottom: 10px; }
				input[type=submit] { padding: 10px 15px; }
			</style>
		</head>
		<body>
			<h2>URL Shortener</h2>
			<form method="post" action="/shorten">
				<input type="url" name="url" placeholder="Enter a URL" required><br>
				<input type="text" name="custom" placeholder="Optional: custom short name"><br>
				<input type="submit" value="Shorten">
			</form>
		</body>
		</html>
	`)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
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
	var shortKey string

	urlMutex.Lock()
	defer urlMutex.Unlock()

	if customKey != "" {
		if _, exists := urlStore[customKey]; exists {
			http.Error(w, "Custom short name is already taken", http.StatusConflict)
			return
		}
		shortKey = customKey
	} else {
		shortKey = generateShortKey()
		// Ensure key is unique
		for _, exists := urlStore[shortKey]; exists; _, exists = urlStore[shortKey] {
			shortKey = generateShortKey()
		}
	}

	urlStore[shortKey] = ShortenedURL{
		ID:          shortKey,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
	}

	shortenedURL := fmt.Sprintf("%s/s/%s", BaseURL, shortKey)
	analyticsURL := fmt.Sprintf("%s/analytics/%s", BaseURL, shortKey)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>URL Shortened</title>
			<style>
				body { font-family: sans-serif; max-width: 600px; margin: 40px auto; padding: 20px; }
				p { margin-bottom: 10px; }
			</style>
		</head>
		<body>
			<h2>URL Shortened!</h2>
			<p>Original URL: %s</p>
			<p>Shortened URL: <a href="%s">%s</a></p>
			<p>Analytics: <a href="%s">%s</a></p>
			<a href="/">Shorten another</a>
		</body>
		</html>
	`, template.HTMLEscapeString(originalURL), template.HTMLEscapeString(shortenedURL), template.HTMLEscapeString(shortenedURL), template.HTMLEscapeString(analyticsURL), template.HTMLEscapeString(analyticsURL))
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := strings.TrimPrefix(r.URL.Path, "/s/")
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	// Check cache first
	cacheMutex.RLock()
	originalURL, found := urlCache[shortKey]
	cacheMutex.RUnlock()

	if !found {
		// If not in cache, check the main store
		urlMutex.RLock()
		urlData, urlFound := urlStore[shortKey]
		urlMutex.RUnlock()

		if !urlFound {
			http.Error(w, "Shortened key not found", http.StatusNotFound)
			return
		}
		originalURL = urlData.OriginalURL

		// Add to cache
		cacheMutex.Lock()
		urlCache[shortKey] = originalURL
		cacheMutex.Unlock()
	}

	// Record the visit for analytics
	visitMutex.Lock()
	visitStore[shortKey] = append(visitStore[shortKey], Visit{
		ShortURLID: shortKey,
		Timestamp:  time.Now(),
		UserAgent:  r.UserAgent(),
	})
	visitMutex.Unlock()

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func handleAnalytics(w http.ResponseWriter, r *http.Request) {
	shortKey := strings.TrimPrefix(r.URL.Path, "/analytics/")
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	urlMutex.RLock()
	urlData, found := urlStore[shortKey]
	urlMutex.RUnlock()
	if !found {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		return
	}

	visitMutex.RLock()
	visits, hasVisits := visitStore[shortKey]
	visitMutex.RUnlock()

	// Aggregation
	dailyCounts := make(map[string]int)
	monthlyCounts := make(map[string]int)
	uaCounts := make(map[string]int)

	if hasVisits {
		for _, visit := range visits {
			day := visit.Timestamp.Format("2006-01-02")
			month := visit.Timestamp.Format("2006-01")
			dailyCounts[day]++
			monthlyCounts[month]++
			uaCounts[visit.UserAgent]++
		}
	}

	// Prepare data for template
	type AnalyticsPageData struct {
		URL             ShortenedURL
		TotalVisits     int
		Visits          []Visit
		DailyCounts     map[string]int
		MonthlyCounts   map[string]int
		UserAgentCounts map[string]int
	}

	data := AnalyticsPageData{
		URL:             urlData,
		TotalVisits:     len(visits),
		Visits:          visits,
		DailyCounts:     dailyCounts,
		MonthlyCounts:   monthlyCounts,
		UserAgentCounts: uaCounts,
	}

	tmpl, err := template.New("analytics").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Analytics</title>
			<style>
				body { font-family: sans-serif; max-width: 800px; margin: 40px auto; padding: 20px; }
				h2, h3 { border-bottom: 1px solid #ccc; padding-bottom: 5px; }
				table { border-collapse: collapse; width: 100%; margin-top: 20px; }
				th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
				th { background-color: #f2f2f2; }
			</style>
		</head>
		<body>
			<h2>Analytics for {{.URL.ID}}</h2>
			<p><strong>Original URL:</strong> {{.URL.OriginalURL}}</p>
			<p><strong>Total Visits:</strong> {{.TotalVisits}}</p>

			<h3>Aggregated by Day</h3>
			<table>
				<tr><th>Day</th><th>Visits</th></tr>
				{{range $day, $count := .DailyCounts}}
				<tr><td>{{$day}}</td><td>{{$count}}</td></tr>
				{{end}}
			</table>

			<h3>Aggregated by Month</h3>
			<table>
				<tr><th>Month</th><th>Visits</th></tr>
				{{range $month, $count := .MonthlyCounts}}
				<tr><td>{{$month}}</td><td>{{$count}}</td></tr>
				{{end}}
			</table>

			<h3>Aggregated by User-Agent</h3>
			<table>
				<tr><th>User-Agent</th><th>Visits</th></tr>
				{{range $ua, $count := .UserAgentCounts}}
				<tr><td>{{$ua}}</td><td>{{$count}}</td></tr>
				{{end}}
			</table>

			<h3>All Visits</h3>
			<table>
				<tr><th>Timestamp</th><th>User-Agent</th></tr>
				{{range .Visits}}
				<tr><td>{{.Timestamp.Format "2006-01-02 15:04:05"}}</td><td>{{.UserAgent}}</td></tr>
				{{end}}
			</table>
			<br>
			<a href="/">Shorten another URL</a>
		</body>
		</html>
	`)
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
	}
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

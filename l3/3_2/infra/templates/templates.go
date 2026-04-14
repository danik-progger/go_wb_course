package templates

import (
	"html/template"
	"io"

	"urlshortener/domain"
)

const FormTemplate = `
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
`

const ShortenResultTemplate = `
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
`

const analyticsTemplateText = `
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
`

var analyticsTemplate = template.Must(template.New("analytics").Parse(analyticsTemplateText))

type AnalyticsData struct {
	URL             domain.ShortenedURL
	TotalVisits     int
	Visits          []domain.Visit
	DailyCounts     map[string]int
	MonthlyCounts   map[string]int
	UserAgentCounts map[string]int
}

func ExecuteAnalytics(w io.Writer, data AnalyticsData) error {
	return analyticsTemplate.Execute(w, data)
}

package parser

import (
	"net/url"
	"strings"
	"testing"
)

func TestResolveURL(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		href     string
		expected string
		wantErr  bool
	}{
		{
			name:     "absolute URL",
			base:     "https://example.com/page/",
			href:     "https://other.com/resource",
			expected: "https://other.com/resource",
			wantErr:  false,
		},
		{
			name:     "relative URL",
			base:     "https://example.com/page/",
			href:     "subpage.html",
			expected: "https://example.com/page/subpage.html",
			wantErr:  false,
		},
		{
			name:     "root-relative URL",
			base:     "https://example.com/page/",
			href:     "/root.html",
			expected: "https://example.com/root.html",
			wantErr:  false,
		},
		{
			name:     "parent directory",
			base:     "https://example.com/a/b/",
			href:     "../c.html",
			expected: "https://example.com/a/c.html",
			wantErr:  false,
		},
		{
			name:    "invalid URL",
			base:    "https://example.com/",
			href:    "://invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseURL, err := url.Parse(tt.base)
			if err != nil {
				t.Fatalf("Failed to parse base URL: %v", err)
			}

			result, err := resolveURL(baseURL, tt.href)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result != tt.expected && !tt.wantErr {
				t.Errorf("resolveURL() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestLinkRewriterToLocalPath(t *testing.T) {
	baseURL, _ := url.Parse("https://example.com/page/")
	rewriter := NewLinkRewriter(baseURL, "/output")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "same domain HTML file",
			input:    "https://example.com/about.html",
			expected: "about.html",
		},
		{
			name:     "same domain directory",
			input:    "https://example.com/contact/",
			expected: "contact/index.html",
		},
		{
			name:     "same domain root",
			input:    "https://example.com/",
			expected: "index.html",
		},
		{
			name:     "different domain",
			input:    "https://other.com/resource",
			expected: "",
		},
		{
			name:     "different domain with path",
			input:    "https://cdn.example.com/style.css",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rewriter.toLocalPath(tt.input)
			if result != tt.expected {
				t.Errorf("toLocalPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractAndRewriteLinks(t *testing.T) {
	baseURL, _ := url.Parse("https://example.com/")
	rewriter := NewLinkRewriter(baseURL, "/output")

	htmlContent := `
<html>
<body>
	<a href="/about.html">About</a>
	<a href="https://other.com/link">External</a>
	<img src="/images/logo.png">
	<script src="/js/app.js"></script>
</body>
</html>
`

	links, content, err := rewriter.ExtractAndRewriteLinks(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("extractAndRewriteLinks() error = %v", err)
	}

	// Check links were extracted
	expectedLinkCount := 4 // about.html, other.com/link, logo.png, app.js
	if len(links) != expectedLinkCount {
		t.Errorf("Expected %d links, got %d", expectedLinkCount, len(links))
	}

	// Check content was rewritten
	rewrittenHTML := content.String()

	// About link should be rewritten to relative path
	if !strings.Contains(rewrittenHTML, `href="about.html"`) {
		t.Error("About link was not rewritten to local path")
	}

	// External link should remain unchanged
	if !strings.Contains(rewrittenHTML, `href="https://other.com/link"`) {
		t.Error("External link was unexpectedly rewritten")
	}

	// Image src should be rewritten
	if !strings.Contains(rewrittenHTML, `src="images/logo.png"`) {
		t.Error("Image src was not rewritten to local path")
	}

	// Script src should be rewritten
	if !strings.Contains(rewrittenHTML, `src="js/app.js"`) {
		t.Error("Script src was not rewritten to local path")
	}
}

package logger

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := LoggingMiddleware(mockHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	logOutput := buf.String()

	if !strings.Contains(logOutput, "🚀 Got request: GET /test") {
		t.Errorf("expected log to contain request info, but it didn't. Log: %s", logOutput)
	}

	if !strings.Contains(logOutput, "🏁 Finished in:") {
		t.Errorf("expected log to contain finish info, but it didn't. Log: %s", logOutput)
	}
}

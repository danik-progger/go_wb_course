package main

import (
	"fmt"
	nethttp "net/http"

	"urlshortener/app/service"
	"urlshortener/infra/db"
	apphttp "urlshortener/infra/http"
)

const (
	BaseURL = "http://localhost:3030"
	Port    = ":3030"
)

func main() {
	urlStore := db.NewURLMemoryStore()
	visitStore := db.NewVisitMemoryStore()
	cache := db.NewMemoryCache()

	svc := service.NewURLService(urlStore, visitStore, cache)
	handler := apphttp.NewHandler(svc, BaseURL)

	nethttp.HandleFunc("/", handler.HandleForm)
	nethttp.HandleFunc("/shorten", handler.HandleShorten)
	nethttp.HandleFunc("/s/", handler.HandleRedirect)
	nethttp.HandleFunc("/analytics/", handler.HandleAnalytics)

	fmt.Println("URL Shortener is running on", Port)
	nethttp.ListenAndServe(Port, nil)
}

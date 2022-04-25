package main

import "net/http"

func setupRoutes() {
	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/readiness", readinessCheckHandler)
	http.HandleFunc("/iplocation", locationHandler)
}

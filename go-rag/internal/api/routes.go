package api

import "net/http"

func RegisterRoutes(h *Handler) {

	http.HandleFunc("/ask", h.Ask)
	http.HandleFunc("/ingest", h.Ingest)
	http.HandleFunc("/health", h.Health)
}
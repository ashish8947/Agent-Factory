package api

import (
	"log"
	"net/http"
)

func StartServer(h *Handler) {

	RegisterRoutes(h)

	log.Println("🚀 RAG API running on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

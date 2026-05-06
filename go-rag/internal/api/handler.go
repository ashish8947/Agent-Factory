package api

import (
	"encoding/json"
	"net/http"

	"github.com/ashish8947/Agent-Factory/go-rag/internal/rag"
)

// -------------------------
// HANDLER STRUCT
// -------------------------

type Handler struct {
	Pipeline *rag.Pipeline
	Ingestor *rag.Ingestor
}

func NewHandler(p *rag.Pipeline, i *rag.Ingestor) *Handler {
	return &Handler{
		Pipeline: p,
		Ingestor: i,
	}
}

// -------------------------
// REQUEST TYPES
// -------------------------

type AskRequest struct {
	Question string `json:"question"`
}

type IngestRequest struct {
	DocID string `json:"doc_id"`
	Text  string `json:"text"`
}

// -------------------------
// ASK (RAG QUERY)
// -------------------------

func (h *Handler) Ask(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var req AskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	answer, err := h.Pipeline.Ask(r.Context(), req.Question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"answer": answer,
	})
}

// -------------------------
// INGEST (ADD DOCUMENTS)
// -------------------------

func (h *Handler) Ingest(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var req IngestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.Ingestor.Ingest(
		r.Context(),
		req.DocID,
		req.Text,
		200, // chunk size
		30,  // overlap
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ingested",
	})
}

// -------------------------
// HEALTH CHECK
// -------------------------

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

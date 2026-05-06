package api

import (
	"encoding/json"
	"net/http"

	"github.com/ashish8947/Agent-Factory/go-rag/internal/rag"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

type AskRequest struct {
	Question string `json:"question"`
}

func (s *Server) Start() {
	s.routes()
	http.ListenAndServe(":8080", nil)
}

func (s *Server) handleAsk(w http.ResponseWriter, r *http.Request) {
	var req AskRequest
	json.NewDecoder(r.Body).Decode(&req)

	answer := rag.Ask(req.Question)

	json.NewEncoder(w).Encode(map[string]string{
		"answer": answer,
	})
}

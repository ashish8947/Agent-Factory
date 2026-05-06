package api

import "net/http"

func (s *Server) routes() {
	http.HandleFunc("/ask", s.handleAsk)
}

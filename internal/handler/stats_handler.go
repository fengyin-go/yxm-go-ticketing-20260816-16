package handler

import (
	"net/http"

	"ticketing/pkg/httpx"
)

// registerStatsRoutes 注册统计相关路由。
func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats", s.stats)
	mux.HandleFunc("GET /api/stats/agents", s.agentWorkloads)
}

func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.Stats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) agentWorkloads(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.AgentWorkloads()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

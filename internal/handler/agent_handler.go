package handler

import (
	"net/http"

	"ticketing/internal/model"
	"ticketing/pkg/httpx"
)

// registerAgentRoutes 注册客服相关路由。
func (s *Server) registerAgentRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/agents", s.createAgent)
	mux.HandleFunc("GET /api/agents", s.listAgents)
	mux.HandleFunc("GET /api/agents/{id}", s.getAgent)
	mux.HandleFunc("PUT /api/agents/{id}", s.updateAgent)
	mux.HandleFunc("PATCH /api/agents/{id}/status", s.setAgentStatus)
	mux.HandleFunc("DELETE /api/agents/{id}", s.deleteAgent)
}

type createAgentRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Department string `json:"department"`
}

func (s *Server) createAgent(w http.ResponseWriter, r *http.Request) {
	var req createAgentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateAgent(model.Agent{
		Name:       req.Name,
		Email:      req.Email,
		Department: req.Department,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listAgents(w http.ResponseWriter, r *http.Request) {
	list, err := s.svc.ListAgents()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}

func (s *Server) getAgent(w http.ResponseWriter, r *http.Request) {
	a, err := s.svc.GetAgent(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

type updateAgentRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Department string `json:"department"`
	Status     string `json:"status"`
}

func (s *Server) updateAgent(w http.ResponseWriter, r *http.Request) {
	var req updateAgentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.UpdateAgent(r.PathValue("id"), model.Agent{
		Name:       req.Name,
		Email:      req.Email,
		Department: req.Department,
		Status:     req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

type setAgentStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) setAgentStatus(w http.ResponseWriter, r *http.Request) {
	var req setAgentStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.SetAgentStatus(r.PathValue("id"), req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteAgent(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteAgent(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}

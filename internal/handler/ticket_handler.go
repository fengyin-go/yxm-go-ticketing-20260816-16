package handler

import (
	"net/http"

	"ticketing/internal/model"
	"ticketing/pkg/httpx"
)

// registerTicketRoutes 注册工单相关路由。
func (s *Server) registerTicketRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/tickets", s.createTicket)
	mux.HandleFunc("GET /api/tickets", s.listTickets)
	mux.HandleFunc("GET /api/tickets/{id}", s.getTicket)
	mux.HandleFunc("PUT /api/tickets/{id}", s.updateTicket)
	mux.HandleFunc("DELETE /api/tickets/{id}", s.deleteTicket)
	mux.HandleFunc("PATCH /api/tickets/{id}/assign", s.assignTicket)
	mux.HandleFunc("PATCH /api/tickets/{id}/status", s.transitionTicket)
}

type createTicketRequest struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	CategoryID     string `json:"category_id"`
	Priority       string `json:"priority"`
	RequesterName  string `json:"requester_name"`
	RequesterEmail string `json:"requester_email"`
}

func (s *Server) createTicket(w http.ResponseWriter, r *http.Request) {
	var req createTicketRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateTicket(model.Ticket{
		Title:          req.Title,
		Description:    req.Description,
		CategoryID:     req.CategoryID,
		Priority:       req.Priority,
		RequesterName:  req.RequesterName,
		RequesterEmail: req.RequesterEmail,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

// listTickets 工单列表：GET /api/tickets?status=&priority=&category_id=&assignee_id=&keyword=&page=&size=
func (s *Server) listTickets(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TicketFilter{
		Status:     r.URL.Query().Get("status"),
		Priority:   r.URL.Query().Get("priority"),
		CategoryID: r.URL.Query().Get("category_id"),
		AssigneeID: r.URL.Query().Get("assignee_id"),
		Keyword:    r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListTickets(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTicket(w http.ResponseWriter, r *http.Request) {
	t, err := s.svc.GetTicket(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type updateTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	CategoryID  string `json:"category_id"`
}

func (s *Server) updateTicket(w http.ResponseWriter, r *http.Request) {
	var req updateTicketRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateTicket(r.PathValue("id"), model.Ticket{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		CategoryID:  req.CategoryID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteTicket(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteTicket(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}

type assignTicketRequest struct {
	AssigneeID string `json:"assignee_id"`
}

func (s *Server) assignTicket(w http.ResponseWriter, r *http.Request) {
	var req assignTicketRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.AssignTicket(r.PathValue("id"), req.AssigneeID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type transitionTicketRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionTicket(w http.ResponseWriter, r *http.Request) {
	var req transitionTicketRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.TransitionTicket(r.PathValue("id"), req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

package handler

import (
	"net/http"

	"ticketing/internal/model"
	"ticketing/pkg/httpx"
)

// registerCommentRoutes 注册留言相关路由。
func (s *Server) registerCommentRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/tickets/{id}/comments", s.addComment)
	mux.HandleFunc("GET /api/tickets/{id}/comments", s.listComments)
	mux.HandleFunc("DELETE /api/comments/{id}", s.deleteComment)
}

type addCommentRequest struct {
	AgentID string `json:"agent_id"`
	Content string `json:"content"`
}

func (s *Server) addComment(w http.ResponseWriter, r *http.Request) {
	var req addCommentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.AddComment(model.Comment{
		TicketID: r.PathValue("id"),
		AgentID:  req.AgentID,
		Content:  req.Content,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listComments(w http.ResponseWriter, r *http.Request) {
	list, err := s.svc.ListComments(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}

func (s *Server) deleteComment(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteComment(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}

package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ticketing/internal/config"
	"ticketing/internal/model"
	"ticketing/internal/service"
	"ticketing/internal/store"
	"ticketing/pkg/httpx"
	"ticketing/pkg/logger"
)

func newListTestServer() *Server {
	cfg := &config.Config{MaxPageSize: 3}
	log := logger.NewLevel(logger.LevelError)
	svc := service.New(store.NewMemoryStore(), log, cfg)
	return NewServer(svc, log, cfg)
}

func seedListTestTickets(t *testing.T, s *Server) {
	t.Helper()
	for i := 1; i <= 5; i++ {
		desc := "普通描述"
		if i == 2 {
			desc = "打印机卡纸，需要现场排查"
		}
		_, err := s.svc.CreateTicket(model.Ticket{
			Title:         "工单标题",
			Description:   desc,
			Priority:      model.PriorityHigh,
			RequesterName: "测试用户",
		})
		if err != nil {
			t.Fatalf("创建工单失败: %v", err)
		}
	}
}

type listResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    httpx.PageResult `json:"data"`
}

func doList(t *testing.T, s *Server, query string) listResponse {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/tickets"+query, nil)
	rr := httptest.NewRecorder()
	s.Routes().ServeHTTP(rr, req)
	var resp listResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("响应解析失败: %v\n%s", err, rr.Body.String())
	}
	return resp
}

func TestListTicketsKeywordMatchesDescription(t *testing.T) {
	s := newListTestServer()
	seedListTestTickets(t, s)

	resp := doList(t, s, "?keyword=打印机")
	if resp.Data.Pagination.Total != 1 {
		t.Fatalf("期望命中 1 条，得到 %d", resp.Data.Pagination.Total)
	}
}

func TestListTicketsPaginationSizeNormalized(t *testing.T) {
	s := newListTestServer()
	seedListTestTickets(t, s)

	resp := doList(t, s, "")
	if resp.Data.Pagination.Size != 3 {
		t.Fatalf("省略 size 时期望归一化为 3，得到 %d", resp.Data.Pagination.Size)
	}
	if got := len(resp.Data.Items.([]interface{})); got != 3 {
		t.Fatalf("期望返回 3 条，得到 %d", got)
	}

	resp = doList(t, s, "?size=1000")
	if resp.Data.Pagination.Size != 3 {
		t.Fatalf("超大 size 时期望限制为 3，得到 %d", resp.Data.Pagination.Size)
	}
}

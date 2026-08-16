package service

import (
	"errors"
	"testing"

	"ticketing/internal/config"
	"ticketing/internal/model"
	"ticketing/internal/store"
	"ticketing/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func (s *Service) seed(t *testing.T) (categoryID, agentID, ticketID string) {
	t.Helper()
	c, err := s.CreateCategory(model.Category{Name: "故障报修"})
	if err != nil {
		t.Fatalf("创建分类失败: %v", err)
	}
	a, err := s.CreateAgent(model.Agent{Name: "张三", Email: "zhang@ex.com"})
	if err != nil {
		t.Fatalf("创建客服失败: %v", err)
	}
	tk, err := s.CreateTicket(model.Ticket{
		Title:         "网络异常",
		Description:   "无法上网",
		CategoryID:    c.ID,
		Priority:      model.PriorityHigh,
		RequesterName: "李四",
	})
	if err != nil {
		t.Fatalf("创建工单失败: %v", err)
	}
	return c.ID, a.ID, tk.ID
}

func TestService_CreateTicket(t *testing.T) {
	s := newTestService()
	// 分类不存在
	if _, err := s.CreateTicket(model.Ticket{Title: "x", CategoryID: "nope", RequesterName: "a"}); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("期望 ErrNotFound，得到 %v", err)
	}
	// 缺少标题
	if _, err := s.CreateTicket(model.Ticket{RequesterName: "a"}); !model.IsValidationError(err) {
		t.Fatalf("期望校验错误，得到 %v", err)
	}
}

func TestService_TicketLifecycle(t *testing.T) {
	s := newTestService()
	_, aid, tid := s.seed(t)

	// 分配客服：open -> processing
	tk, err := s.AssignTicket(tid, aid)
	if err != nil {
		t.Fatalf("分配失败: %v", err)
	}
	if tk.Status != model.StatusProcessing {
		t.Fatalf("期望 processing，得到 %s", tk.Status)
	}

	// processing -> resolved
	tk, err = s.TransitionTicket(tid, model.StatusResolved)
	if err != nil {
		t.Fatalf("流转失败: %v", err)
	}
	if tk.ResolvedAt == nil {
		t.Fatalf("解决后应记录 ResolvedAt")
	}

	// resolved -> closed
	if _, err = s.TransitionTicket(tid, model.StatusClosed); err != nil {
		t.Fatalf("关闭失败: %v", err)
	}

	// closed -> processing 非法
	if _, err = s.TransitionTicket(tid, model.StatusProcessing); !model.IsValidationError(err) {
		t.Fatalf("期望校验错误，得到 %v", err)
	}
}

func TestService_Comments(t *testing.T) {
	s := newTestService()
	_, aid, tid := s.seed(t)

	c, err := s.AddComment(model.Comment{TicketID: tid, AgentID: aid, Content: "已定位问题"})
	if err != nil {
		t.Fatalf("添加留言失败: %v", err)
	}
	if c.ID == "" {
		t.Fatalf("留言 ID 为空")
	}

	list, err := s.ListComments(tid)
	if err != nil {
		t.Fatalf("列出留言失败: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("期望 1 条留言，得到 %d", len(list))
	}
}

func TestService_StatsAndWorkload(t *testing.T) {
	s := newTestService()
	cid, aid, _ := s.seed(t)

	stats, err := s.Stats()
	if err != nil {
		t.Fatalf("统计失败: %v", err)
	}
	if stats.Total != 1 || stats.Unassigned != 1 {
		t.Fatalf("统计异常: %+v", stats)
	}

	// 再建一个工单并分配，验证工作量统计
	tk2, err := s.CreateTicket(model.Ticket{
		Title:         "第二个工单",
		CategoryID:    cid,
		Priority:      model.PriorityLow,
		RequesterName: "王五",
	})
	if err != nil {
		t.Fatalf("创建工单失败: %v", err)
	}
	if _, err := s.AssignTicket(tk2.ID, aid); err != nil {
		t.Fatalf("分配失败: %v", err)
	}
	workloads, _ := s.AgentWorkloads()
	if len(workloads) == 0 || workloads[0].Total != 1 {
		t.Fatalf("工作量异常: %+v", workloads)
	}
}

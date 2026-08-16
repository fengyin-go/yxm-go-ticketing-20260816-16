package store

import (
	"testing"
	"time"

	"ticketing/internal/model"
)

func TestMemoryStore_Category(t *testing.T) {
	s := NewMemoryStore()
	c := &model.Category{ID: "c1", Name: "故障报修", CreatedAt: time.Now()}
	if err := s.CreateCategory(c); err != nil {
		t.Fatalf("创建分类失败: %v", err)
	}
	if err := s.CreateCategory(&model.Category{ID: "c2", Name: "故障报修"}); err != ErrConflict {
		t.Fatalf("期望 ErrConflict，得到 %v", err)
	}
	got, err := s.GetCategory("c1")
	if err != nil || got.Name != "故障报修" {
		t.Fatalf("查询分类失败: %v, %v", got, err)
	}
	if len(s.ListCategories()) != 1 {
		t.Fatalf("期望 1 个分类")
	}
}

func TestMemoryStore_Agent(t *testing.T) {
	s := NewMemoryStore()
	a := &model.Agent{ID: "a1", Name: "张三", Email: "zhang@ex.com", Status: model.AgentAvailable, CreatedAt: time.Now()}
	if err := s.CreateAgent(a); err != nil {
		t.Fatalf("创建客服失败: %v", err)
	}
	if err := s.CreateAgent(&model.Agent{ID: "a2", Email: "zhang@ex.com"}); err != ErrConflict {
		t.Fatalf("期望 ErrConflict，得到 %v", err)
	}
	got, err := s.GetAgent("a1")
	if err != nil || got.Name != "张三" {
		t.Fatalf("查询客服失败: %v, %v", got, err)
	}
}

func TestMemoryStore_TicketAndComment(t *testing.T) {
	s := NewMemoryStore()
	tk := &model.Ticket{ID: "t1", Title: "网络异常", Status: model.StatusOpen, Priority: model.PriorityHigh, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateTicket(tk); err != nil {
		t.Fatalf("创建工单失败: %v", err)
	}
	if len(s.ListTickets()) != 1 {
		t.Fatalf("期望 1 个工单")
	}

	c := &model.Comment{ID: "cm1", TicketID: "t1", Content: "已排查", CreatedAt: time.Now()}
	if err := s.CreateComment(c); err != nil {
		t.Fatalf("创建留言失败: %v", err)
	}
	if got := len(s.ListComments("t1")); got != 1 {
		t.Fatalf("期望 1 条留言，得到 %d", got)
	}

	tk.Status = model.StatusResolved
	if err := s.UpdateTicket(tk); err != nil {
		t.Fatalf("更新工单失败: %v", err)
	}
	got, _ := s.GetTicket("t1")
	if got.Status != model.StatusResolved {
		t.Fatalf("期望已解决，得到 %s", got.Status)
	}
}

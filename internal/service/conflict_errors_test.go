package service

import (
	"errors"
	"testing"

	"ticketing/internal/model"
	"ticketing/internal/store"
)

func TestConflictErrorsPreserveSentinel(t *testing.T) {
	s := newTestService()

	if _, err := s.CreateCategory(model.Category{Name: "重复分类"}); err != nil {
		t.Fatalf("首次创建分类失败: %v", err)
	}
	if _, err := s.CreateCategory(model.Category{Name: "重复分类"}); !errors.Is(err, store.ErrConflict) {
		t.Fatalf("分类冲突应保留 ErrConflict，得到 %v", err)
	}

	if _, err := s.CreateAgent(model.Agent{Name: "张三", Email: "dup@ex.com"}); err != nil {
		t.Fatalf("首次创建客服失败: %v", err)
	}
	if _, err := s.CreateAgent(model.Agent{Name: "李四", Email: "dup@ex.com"}); !errors.Is(err, store.ErrConflict) {
		t.Fatalf("客服冲突应保留 ErrConflict，得到 %v", err)
	}
}

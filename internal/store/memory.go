package store

import (
	"sync"

	"ticketing/internal/model"
)

// MemoryStore 基于内存的 Store 实现。
type MemoryStore struct {
	mu         sync.RWMutex
	categories map[string]*model.Category
	agents    map[string]*model.Agent
	tickets   map[string]*model.Ticket
	comments  map[string]*model.Comment
}

func cloneCategory(src *model.Category) *model.Category {
	if src == nil {
		return nil
	}
	c := *src
	return &c
}

func cloneAgent(src *model.Agent) *model.Agent {
	if src == nil {
		return nil
	}
	a := *src
	return &a
}

func cloneTicket(src *model.Ticket) *model.Ticket {
	if src == nil {
		return nil
	}
	t := *src
	if src.ResolvedAt != nil {
		resolvedAt := *src.ResolvedAt
		t.ResolvedAt = &resolvedAt
	}
	return &t
}

func cloneComment(src *model.Comment) *model.Comment {
	if src == nil {
		return nil
	}
	c := *src
	return &c
}

// NewMemoryStore 创建空的内存存储。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		categories: make(map[string]*model.Category),
		agents:    make(map[string]*model.Agent),
		tickets:   make(map[string]*model.Ticket),
		comments:  make(map[string]*model.Comment),
	}
}

var _ Store = (*MemoryStore)(nil)

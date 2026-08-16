package store

import (
	"fmt"

	"ticketing/internal/model"
)

// CreateAgent 新增客服，邮箱重复时返回 ErrConflict。
func (s *MemoryStore) CreateAgent(a *model.Agent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.agents {
		if exist.Email == a.Email {
			return fmt.Errorf("agent conflict: %v", ErrConflict)
		}
	}
	s.agents[a.ID] = a
	return nil
}

// GetAgent 按 ID 查询客服。
func (s *MemoryStore) GetAgent(id string) (*model.Agent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.agents[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

// ListAgents 返回全部客服。
func (s *MemoryStore) ListAgents() []*model.Agent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Agent, 0, len(s.agents))
	for _, a := range s.agents {
		list = append(list, a)
	}
	return list
}

// UpdateAgent 覆盖保存客服。
func (s *MemoryStore) UpdateAgent(a *model.Agent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.agents[a.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.agents {
		if exist.ID != a.ID && exist.Email == a.Email {
			return fmt.Errorf("agent conflict: %v", ErrConflict)
		}
	}
	s.agents[a.ID] = a
	return nil
}

// DeleteAgent 按 ID 删除客服。
func (s *MemoryStore) DeleteAgent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.agents[id]; !ok {
		return ErrNotFound
	}
	delete(s.agents, id)
	return nil
}

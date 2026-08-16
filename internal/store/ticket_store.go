package store

import "ticketing/internal/model"

// CreateTicket 新增工单。
func (s *MemoryStore) CreateTicket(t *model.Ticket) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tickets[t.ID] = t
	return nil
}

// GetTicket 按 ID 查询工单。
func (s *MemoryStore) GetTicket(id string) (*model.Ticket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tickets[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cloneTicket(t), nil
}

// ListTickets 返回全部工单。
func (s *MemoryStore) ListTickets() []*model.Ticket {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Ticket, 0, len(s.tickets))
	for _, t := range s.tickets {
		list = append(list, cloneTicket(t))
	}
	return list
}

// UpdateTicket 覆盖保存工单。
func (s *MemoryStore) UpdateTicket(t *model.Ticket) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tickets[t.ID]; !ok {
		return ErrNotFound
	}
	s.tickets[t.ID] = cloneTicket(t)
	return nil
}

// DeleteTicket 按 ID 删除工单。
func (s *MemoryStore) DeleteTicket(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tickets[id]; !ok {
		return ErrNotFound
	}
	delete(s.tickets, id)
	return nil
}

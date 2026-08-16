package store

import "ticketing/internal/model"

// CreateComment 新增留言。
func (s *MemoryStore) CreateComment(c *model.Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.comments[c.ID] = c
	return nil
}

// GetComment 按 ID 查询留言。
func (s *MemoryStore) GetComment(id string) (*model.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.comments[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

// ListComments 返回指定工单的全部留言。
func (s *MemoryStore) ListComments(ticketID string) []*model.Comment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Comment, 0)
	for _, c := range s.comments {
		if c.TicketID == ticketID {
			list = append(list, c)
		}
	}
	return list
}

// DeleteComment 按 ID 删除留言。
func (s *MemoryStore) DeleteComment(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.comments[id]; !ok {
		return ErrNotFound
	}
	delete(s.comments, id)
	return nil
}

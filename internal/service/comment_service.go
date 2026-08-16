package service

import (
	"sort"
	"time"

	"ticketing/internal/model"
	"ticketing/pkg/idgen"
)

// AddComment 为工单添加留言（处理记录）。
func (s *Service) AddComment(input model.Comment) (*model.Comment, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetTicket(input.TicketID); err != nil {
		return nil, err
	}
	if input.AgentID != "" {
		if _, err := s.store.GetAgent(input.AgentID); err != nil {
			return nil, err
		}
	}
	c := &model.Comment{
		ID:        idgen.Hex(),
		TicketID:  input.TicketID,
		AgentID:   input.AgentID,
		Content:   input.Content,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateComment(c); err != nil {
		return nil, err
	}
	s.log.Infof("工单 %s 新增留言", input.TicketID)
	return c, nil
}

// ListComments 列出某工单的全部留言，按时间正序。
func (s *Service) ListComments(ticketID string) ([]*model.Comment, error) {
	if _, err := s.store.GetTicket(ticketID); err != nil {
		return nil, err
	}
	list := s.store.ListComments(ticketID)
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})
	return list, nil
}

// DeleteComment 删除留言。
func (s *Service) DeleteComment(id string) error {
	if err := s.store.DeleteComment(id); err != nil {
		return err
	}
	s.log.Infof("删除留言 %s", id)
	return nil
}

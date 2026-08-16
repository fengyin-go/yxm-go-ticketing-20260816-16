package service

import (
	"sort"
	"time"

	"ticketing/internal/model"
	"ticketing/pkg/idgen"
)

// priorityWeight 用于排序的优先级权重，数值越大越靠前。
var priorityWeight = map[string]int{
	model.PriorityUrgent: 4,
	model.PriorityHigh:   3,
	model.PriorityMedium: 2,
	model.PriorityLow:    1,
}

// CreateTicket 创建工单。
func (s *Service) CreateTicket(input model.Ticket) (*model.Ticket, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if input.CategoryID != "" {
		if _, err := s.store.GetCategory(input.CategoryID); err != nil {
			return nil, err
		}
	}
	t := &model.Ticket{
		ID:             idgen.Hex(),
		Title:          input.Title,
		Description:    input.Description,
		CategoryID:     input.CategoryID,
		Priority:       input.Priority,
		Status:         model.StatusOpen,
		RequesterName:  input.RequesterName,
		RequesterEmail: input.RequesterEmail,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.store.CreateTicket(t); err != nil {
		return nil, err
	}
	s.log.Infof("创建工单 %s", t.Title)
	return t, nil
}

// GetTicket 按 ID 查询工单。
func (s *Service) GetTicket(id string) (*model.Ticket, error) {
	return s.store.GetTicket(id)
}

// ListTickets 列出工单，支持筛选、按优先级与时间排序、分页。
func (s *Service) ListTickets(filter model.TicketFilter, page, size int) ([]*model.Ticket, int, error) {
	all := s.store.ListTickets()
	matched := make([]*model.Ticket, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		pi, pj := priorityWeight[matched[i].Priority], priorityWeight[matched[j].Priority]
		if pi != pj {
			return pi > pj
		}
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Ticket{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateTicket 更新工单的标题、描述、优先级与分类。
func (s *Service) UpdateTicket(id string, input model.Ticket) (*model.Ticket, error) {
	exist, err := s.store.GetTicket(id)
	if err != nil {
		return nil, err
	}
	exist.Title = input.Title
	exist.Description = input.Description
	exist.Priority = input.Priority
	exist.CategoryID = input.CategoryID
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateTicket(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

// AssignTicket 将工单分配给客服。
func (s *Service) AssignTicket(id, assigneeID string) (*model.Ticket, error) {
	exist, err := s.store.GetTicket(id)
	if err != nil {
		return nil, err
	}
	if assigneeID != "" {
		if _, err := s.store.GetAgent(assigneeID); err != nil {
			return nil, err
		}
	}
	exist.AssigneeID = assigneeID
	if assigneeID != "" {
		if exist.Status == model.StatusOpen {
			exist.Status = model.StatusProcessing
		}
	} else if exist.Status == model.StatusProcessing {
		exist.Status = model.StatusOpen
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateTicket(exist); err != nil {
		return nil, err
	}
	s.log.Infof("工单 %s 分配给 %s", id, assigneeID)
	return exist, nil
}

// TransitionTicket 变更工单状态，校验状态机合法性。
func (s *Service) TransitionTicket(id, status string) (*model.Ticket, error) {
	exist, err := s.store.GetTicket(id)
	if err != nil {
		return nil, err
	}
	if !model.ValidStatus(status) {
		return nil, model.NewValidationError("status", "工单状态不合法")
	}
	if !model.CanTransition(exist.Status, status) {
		return nil, model.NewValidationError("status",
			"不允许从 "+exist.Status+" 流转到 "+status)
	}
	exist.Status = status
	now := time.Now()
	exist.UpdatedAt = now
	if status == model.StatusResolved {
		exist.ResolvedAt = &now
	}
	if err := s.store.UpdateTicket(exist); err != nil {
		return nil, err
	}
	s.log.Infof("工单 %s 状态变更为 %s", id, status)
	return exist, nil
}

// DeleteTicket 删除工单。
func (s *Service) DeleteTicket(id string) error {
	if err := s.store.DeleteTicket(id); err != nil {
		return err
	}
	for _, c := range s.store.ListComments(id) {
		_ = s.store.DeleteComment(c.ID)
	}
	s.log.Infof("删除工单 %s", id)
	return nil
}

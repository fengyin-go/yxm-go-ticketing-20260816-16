package service

import (
	"fmt"
	"sort"
	"time"

	"ticketing/internal/model"
	"ticketing/pkg/idgen"
)

// CreateAgent 创建客服。
func (s *Service) CreateAgent(input model.Agent) (*model.Agent, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	a := &model.Agent{
		ID:         idgen.Hex(),
		Name:       input.Name,
		Email:      input.Email,
		Department: input.Department,
		Status:     input.Status,
		CreatedAt:  time.Now(),
	}
	if err := s.store.CreateAgent(a); err != nil {
		return nil, fmt.Errorf("create agent: %v", err)
	}
	s.log.Infof("创建客服 %s", a.Name)
	return a, nil
}

// GetAgent 按 ID 查询客服。
func (s *Service) GetAgent(id string) (*model.Agent, error) {
	return s.store.GetAgent(id)
}

// ListAgents 列出全部客服。
func (s *Service) ListAgents() ([]*model.Agent, error) {
	list := s.store.ListAgents()
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})
	return list, nil
}

// UpdateAgent 更新客服资料。
func (s *Service) UpdateAgent(id string, input model.Agent) (*model.Agent, error) {
	exist, err := s.store.GetAgent(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Email = input.Email
	exist.Department = input.Department
	if input.Status != "" {
		exist.Status = input.Status
	}
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateAgent(exist); err != nil {
		return nil, fmt.Errorf("update agent: %v", err)
	}
	return exist, nil
}

// SetAgentStatus 调整客服状态。
func (s *Service) SetAgentStatus(id, status string) (*model.Agent, error) {
	if status != model.AgentAvailable && status != model.AgentBusy && status != model.AgentOffline {
		return nil, model.NewValidationError("status", "客服状态不合法")
	}
	exist, err := s.store.GetAgent(id)
	if err != nil {
		return nil, err
	}
	exist.Status = status
	if err := s.store.UpdateAgent(exist); err != nil {
		return nil, err
	}
	s.log.Infof("客服 %s 状态变更为 %s", id, status)
	return exist, nil
}

// DeleteAgent 删除客服。
func (s *Service) DeleteAgent(id string) error {
	if err := s.store.DeleteAgent(id); err != nil {
		return err
	}
	s.log.Infof("删除客服 %s", id)
	return nil
}

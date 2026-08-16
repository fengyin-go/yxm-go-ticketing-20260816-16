package service

import (
	"sort"

	"ticketing/internal/model"
)

// TicketStats 工单全局统计。
type TicketStats struct {
	Total        int            `json:"total"`
	ByStatus     map[string]int `json:"by_status"`
	ByPriority   map[string]int `json:"by_priority"`
	Unassigned   int            `json:"unassigned"`
}

// Stats 汇总工单统计。
func (s *Service) Stats() (*TicketStats, error) {
	stats := &TicketStats{
		ByStatus:   make(map[string]int),
		ByPriority: make(map[string]int),
	}
	for _, t := range s.store.ListTickets() {
		stats.Total++
		stats.ByStatus[t.Status]++
		stats.ByPriority[t.Priority]++
		if t.AssigneeID == "" {
			stats.Unassigned++
		}
	}
	return stats, nil
}

// AgentWorkload 客服工作量统计。
type AgentWorkload struct {
	AgentID   string `json:"agent_id"`
	AgentName string `json:"agent_name"`
	Total     int    `json:"total"`
	Open      int    `json:"open"`
	Resolved  int    `json:"resolved"`
}

// AgentWorkloads 统计各客服的工单工作量。
func (s *Service) AgentWorkloads() ([]*AgentWorkload, error) {
	agents := s.store.ListAgents()
	nameByID := make(map[string]string, len(agents))
	workload := make(map[string]*AgentWorkload, len(agents))
	for _, a := range agents {
		nameByID[a.ID] = a.Name
		workload[a.ID] = &AgentWorkload{AgentID: a.ID, AgentName: a.Name}
	}
	for _, t := range s.store.ListTickets() {
		if t.AssigneeID == "" {
			continue
		}
		wl, ok := workload[t.AssigneeID]
		if !ok {
			wl = &AgentWorkload{AgentID: t.AssigneeID, AgentName: nameByID[t.AssigneeID]}
			workload[t.AssigneeID] = wl
		}
		wl.Total++
		if t.Status == model.StatusOpen {
			wl.Open++
		}
		if t.Status == model.StatusResolved {
			wl.Resolved++
		}
	}
	result := make([]*AgentWorkload, 0, len(workload))
	for _, wl := range workload {
		result = append(result, wl)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Total > result[j].Total
	})
	return result, nil
}

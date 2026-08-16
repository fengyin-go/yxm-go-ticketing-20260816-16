package model

import (
	"strings"
	"time"
)

// 工单状态常量。
const (
	StatusOpen       = "open"       // 待处理
	StatusProcessing = "processing" // 处理中
	StatusResolved   = "resolved"   // 已解决
	StatusClosed     = "closed"     // 已关闭
)

// 工单优先级常量。
const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
	PriorityUrgent = "urgent"
)

// ticketTransitions 定义工单状态的合法流转关系。
var ticketTransitions = map[string]map[string]bool{
	StatusOpen:       {StatusProcessing: true, StatusResolved: true, StatusClosed: true},
	StatusProcessing: {StatusResolved: true, StatusClosed: true},
	StatusResolved:   {StatusClosed: true, StatusOpen: true},
	StatusClosed:     {StatusOpen: true},
}

// Ticket 工单实体。
type Ticket struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	CategoryID     string     `json:"category_id"`
	Priority       string     `json:"priority"`
	Status         string     `json:"status"`
	AssigneeID     string     `json:"assignee_id,omitempty"` // 负责客服
	RequesterName  string     `json:"requester_name"`
	RequesterEmail string     `json:"requester_email"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
}

// Validate 规范化并校验工单字段。
func (t *Ticket) Validate() error {
	t.Title = strings.TrimSpace(t.Title)
	t.Description = strings.TrimSpace(t.Description)
	t.RequesterName = strings.TrimSpace(t.RequesterName)
	t.RequesterEmail = strings.TrimSpace(t.RequesterEmail)
	if t.Title == "" {
		return NewValidationError("title", "工单标题不能为空")
	}
	if t.RequesterName == "" {
		return NewValidationError("requester_name", "请求人姓名不能为空")
	}
	if t.Priority == "" {
		t.Priority = PriorityMedium
	}
	if !validPriority(t.Priority) {
		return NewValidationError("priority", "优先级不合法")
	}
	if t.Status == "" {
		t.Status = StatusOpen
	}
	if !validStatus(t.Status) {
		return NewValidationError("status", "工单状态不合法")
	}
	return nil
}

// CanTransition 判断从当前状态能否流转到目标状态。
func CanTransition(from, to string) bool {
	if m, ok := ticketTransitions[from]; ok {
		return m[to]
	}
	return false
}

// validStatus 校验状态是否合法。
func validStatus(s string) bool {
	switch s {
	case StatusOpen, StatusProcessing, StatusResolved, StatusClosed:
		return true
	default:
		return false
	}
}

// ValidStatus 对外暴露的状态校验。
func ValidStatus(s string) bool { return validStatus(s) }

// validPriority 校验优先级是否合法。
func validPriority(p string) bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent:
		return true
	default:
		return false
	}
}

// ValidPriority 对外暴露的优先级校验。
func ValidPriority(p string) bool { return validPriority(p) }

// TicketFilter 工单列表筛选条件。
type TicketFilter struct {
	Status     string
	Priority   string
	CategoryID string
	AssigneeID string
	Keyword    string // 匹配标题
}

// Match 判断工单是否命中筛选条件。
func (f TicketFilter) Match(t *Ticket) bool {
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.Priority != "" && t.Priority != f.Priority {
		return false
	}
	if f.CategoryID != "" && t.CategoryID != f.CategoryID {
		return false
	}
	if f.AssigneeID != "" && t.AssigneeID != f.AssigneeID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.Title), k) {
			return false
		}
	}
	return true
}

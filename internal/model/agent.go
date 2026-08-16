package model

import (
	"strings"
	"time"
)

// 客服状态常量。
const (
	AgentAvailable = "available" // 空闲
	AgentBusy      = "busy"      // 忙碌
	AgentOffline   = "offline"   // 离线
)

// Agent 客服人员。
type Agent struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Department string    `json:"department"`
	Status     string    `json:"status"` // available / busy / offline
	CreatedAt  time.Time `json:"created_at"`
}

// Validate 规范化并校验客服字段。
func (a *Agent) Validate() error {
	a.Name = strings.TrimSpace(a.Name)
	a.Email = strings.TrimSpace(a.Email)
	a.Department = strings.TrimSpace(a.Department)
	if a.Name == "" {
		return NewValidationError("name", "客服姓名不能为空")
	}
	if a.Email == "" || !strings.Contains(a.Email, "@") {
		return NewValidationError("email", "邮箱格式不正确")
	}
	if a.Status == "" {
		a.Status = AgentAvailable
	}
	if a.Status != AgentAvailable && a.Status != AgentBusy && a.Status != AgentOffline {
		return NewValidationError("status", "客服状态不合法")
	}
	return nil
}

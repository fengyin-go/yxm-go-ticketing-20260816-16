package model

import (
	"strings"
	"time"
)

// Comment 工单留言（由客服添加的处理记录）。
type Comment struct {
	ID        string    `json:"id"`
	TicketID  string    `json:"ticket_id"`
	AgentID   string    `json:"agent_id,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 规范化并校验留言字段。
func (c *Comment) Validate() error {
	c.TicketID = strings.TrimSpace(c.TicketID)
	c.Content = strings.TrimSpace(c.Content)
	if c.TicketID == "" {
		return NewValidationError("ticket_id", "工单不能为空")
	}
	if c.Content == "" {
		return NewValidationError("content", "留言内容不能为空")
	}
	return nil
}

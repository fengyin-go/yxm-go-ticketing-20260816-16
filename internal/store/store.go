// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"ticketing/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法。
type Store interface {
	// 分类
	CreateCategory(c *model.Category) error
	GetCategory(id string) (*model.Category, error)
	ListCategories() []*model.Category
	UpdateCategory(c *model.Category) error
	DeleteCategory(id string) error

	// 客服
	CreateAgent(a *model.Agent) error
	GetAgent(id string) (*model.Agent, error)
	ListAgents() []*model.Agent
	UpdateAgent(a *model.Agent) error
	DeleteAgent(id string) error

	// 工单
	CreateTicket(t *model.Ticket) error
	GetTicket(id string) (*model.Ticket, error)
	ListTickets() []*model.Ticket
	UpdateTicket(t *model.Ticket) error
	DeleteTicket(id string) error
	CountTicketsByCategory(categoryID string) int

	// 留言
	CreateComment(c *model.Comment) error
	GetComment(id string) (*model.Comment, error)
	ListComments(ticketID string) []*model.Comment
	DeleteComment(id string) error
}

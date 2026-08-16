// Package service 实现工单系统的业务逻辑层。
package service

import (
	"ticketing/internal/config"
	"ticketing/internal/store"
	"ticketing/pkg/logger"
)

// Service 业务逻辑层入口。
type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

// New 创建业务服务实例。
func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}

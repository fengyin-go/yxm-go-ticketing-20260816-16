package service

import (
	"errors"
	"testing"

	"ticketing/internal/model"
	"ticketing/internal/store"
)

func TestAssignUnassignReopensProcessingTicket(t *testing.T) {
	s := newTestService()
	_, aid, tid := s.seed(t)

	tk, err := s.AssignTicket(tid, aid)
	if err != nil {
		t.Fatalf("分配失败: %v", err)
	}
	if tk.Status != model.StatusProcessing {
		t.Fatalf("期望 processing，得到 %s", tk.Status)
	}

	tk, err = s.AssignTicket(tid, "")
	if err != nil {
		t.Fatalf("取消分配失败: %v", err)
	}
	if tk.Status != model.StatusOpen {
		t.Fatalf("取消分配后期望 open，得到 %s", tk.Status)
	}
}

func TestDeleteAssignedAgentRejected(t *testing.T) {
	s := newTestService()
	_, aid, tid := s.seed(t)

	if _, err := s.AssignTicket(tid, aid); err != nil {
		t.Fatalf("分配失败: %v", err)
	}
	if err := s.DeleteAgent(aid); !errors.Is(err, store.ErrConflict) {
		t.Fatalf("期望 ErrConflict，得到 %v", err)
	}
}

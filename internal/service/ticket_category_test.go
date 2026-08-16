package service

import (
	"errors"
	"testing"

	"ticketing/internal/model"
	"ticketing/internal/store"
)

func TestUpdateTicketRejectsMissingCategory(t *testing.T) {
	s := newTestService()
	_, _, tid := s.seed(t)

	_, err := s.UpdateTicket(tid, model.Ticket{
		Title:       "改标题",
		Description: "改描述",
		Priority:    model.PriorityHigh,
		CategoryID:  "not-exists",
	})
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("期望 ErrNotFound，得到 %v", err)
	}
}

func TestDeleteCategoryRejectsUsedCategory(t *testing.T) {
	s := newTestService()
	cid, _, _ := s.seed(t)

	if err := s.DeleteCategory(cid); !errors.Is(err, store.ErrConflict) {
		t.Fatalf("期望 ErrConflict，得到 %v", err)
	}
}

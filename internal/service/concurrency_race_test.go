package service

import (
	"sync"
	"testing"

	"ticketing/internal/model"
)

func TestConcurrentAgentMutationRace(t *testing.T) {
	s := newTestService()
	a, err := s.CreateAgent(model.Agent{Name: "初始客服", Email: "race-agent@ex.com"})
	if err != nil {
		t.Fatalf("创建客服失败: %v", err)
	}

	stop := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					got, _ := s.GetAgent(a.ID)
					if got != nil {
						_ = got.Name
						_ = got.Email
					}
				}
			}
		}()
	}

	for i := 0; i < 80; i++ {
		if _, err := s.UpdateAgent(a.ID, model.Agent{Name: "并发更新", Email: "race-agent@ex.com"}); err != nil {
			t.Fatalf("更新客服失败: %v", err)
		}
	}
	close(stop)
	wg.Wait()
}

func TestConcurrentTicketMutationRace(t *testing.T) {
	s := newTestService()
	tk, err := s.CreateTicket(model.Ticket{
		Title:         "并发工单",
		Priority:      model.PriorityHigh,
		RequesterName: "并发用户",
	})
	if err != nil {
		t.Fatalf("创建工单失败: %v", err)
	}

	stop := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					got, _ := s.GetTicket(tk.ID)
					if got != nil {
						_ = got.Title
						_ = got.Priority
					}
				}
			}
		}()
	}

	for i := 0; i < 80; i++ {
		if _, err := s.UpdateTicket(tk.ID, model.Ticket{
			Title:    "并发更新标题",
			Priority: model.PriorityMedium,
		}); err != nil {
			t.Fatalf("更新工单失败: %v", err)
		}
	}
	close(stop)
	wg.Wait()
}

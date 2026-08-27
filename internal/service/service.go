package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/tide-archive/internal/store"
	"sync"
	"time"
)

type Service struct {
	repo *store.Repository
	jobs chan string
	stop chan struct{}
	once sync.Once
}

func New(repo *store.Repository) *Service {
	s := &Service{repo: repo, jobs: make(chan string, 32), stop: make(chan struct{})}
	go s.runExporter()
	return s
}
func (s *Service) Close() { s.once.Do(func() { close(s.stop) }) }
func (s *Service) CreateSession(ctx context.Context, id, site string) error {
	if err := ready(ctx, id, site); err != nil {
		return err
	}
	if err := s.repo.CreateSession(ctx, id, site, time.Now().UTC().Format(time.RFC3339)); err != nil {
		return fmt.Errorf("create acoustic session: %w", err)
	}
	return s.repo.AppendAudit(ctx, id, "session created")
}
func (s *Service) StartSession(ctx context.Context, id string) error {
	if err := s.repo.SetState(ctx, id, "recording"); err != nil {
		return fmt.Errorf("start archive session: %w", err)
	}
	return s.repo.AppendAudit(ctx, id, "recording started")
}
func (s *Service) RegisterFragment(ctx context.Context, id, sessionID, label string, score float64) error {
	if err := ready(ctx, id, sessionID); err != nil {
		return err
	}
	if err := s.repo.AddFragment(ctx, id, sessionID, label, score); err != nil {
		return fmt.Errorf("register sound fragment: %w", err)
	}
	return s.repo.AppendAudit(ctx, sessionID, "fragment registered")
}
func (s *Service) ReviewFragment(ctx context.Context, sessionID, fragmentID string) error {
	if err := s.repo.ReviewFragment(ctx, fragmentID); err != nil {
		return fmt.Errorf("review candidate: %w", err)
	}
	return s.repo.AppendAudit(ctx, sessionID, "fragment reviewed")
}
func (s *Service) CloseSession(ctx context.Context, id string) error {
	open, err := s.repo.CountUnreviewed(ctx, id)
	if err != nil {
		return fmt.Errorf("count pending reviews: %w", err)
	}
	if open != 0 {
		return fmt.Errorf("session %s still has %d unreviewed fragments", id, open)
	}
	if err := s.repo.SetState(ctx, id, "closed"); err != nil {
		return err
	}
	return s.repo.AppendAudit(ctx, id, "session closed")
}
func (s *Service) Archive(ctx context.Context, id string) error {
	if err := s.repo.SetState(ctx, id, "archived"); err != nil {
		return err
	}
	select {
	case s.jobs <- id:
		return s.repo.AppendAudit(ctx, id, "archive export queued")
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (s *Service) runExporter() {
	for {
		select {
		case <-s.stop:
			return
		case id := <-s.jobs:
			_ = s.repo.AppendAudit(context.Background(), id, "archive export completed")
		}
	}
}
func ready(ctx context.Context, values ...string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	for _, value := range values {
		if value == "" {
			return fmt.Errorf("required archive value is empty")
		}
	}
	return nil
}

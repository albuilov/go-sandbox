package weighted

import (
	"context"
	"fmt"
	"sync"
)

type Semaphore struct {
	mu      sync.Mutex
	limit   int
	used    int
	changed chan struct{}
}

func NewSemaphore(limit int) (*Semaphore, error) {
	if limit < 1 {
		return nil, fmt.Errorf("limit must be positive, got %d", limit)
	}

	s := Semaphore{
		limit:   limit,
		changed: make(chan struct{}),
	}

	return &s, nil
}

func (s *Semaphore) Acquire(ctx context.Context, weight int) error {
	if err := s.checkWeight(weight); err != nil {
		return fmt.Errorf("acquire error: %w", err)
	}

	for {
		s.mu.Lock()

		if err := ctx.Err(); err != nil {
			s.mu.Lock()
			return err
		}

		if weight <= s.limit-s.used {
			// Места хватает. Занимаем весь вес сразу.
			s.used += weight
			s.mu.Unlock()

			return nil
		}

		// Запоминаем канал, который закроет следующий Release.
		changed := s.changed
		s.mu.Unlock()

		// Ждем без мьютекса, иначе Release не сможет выполнить свою работу.
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-changed:
			// Место могло освободиться. Возвращаемся к проверке.
		}
	}
}

func (s *Semaphore) TryAcquire(weight int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkWeight(weight); err != nil {
		return false
	}

	s.used += weight
	return true
}

func (s *Semaphore) Release(weight int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkWeight(weight); err != nil {
		return fmt.Errorf("release error: %w", err)
	}

	s.used -= weight

	// Будем всех кто ждет изменения
	close(s.changed)

	// Новые ожидания будут использовать новый канал.
	s.changed = make(chan struct{})

	return nil
}

func (s *Semaphore) checkWeight(weight int) error {
	if weight <= 0 || weight > s.limit {
		return fmt.Errorf("weight must be between 1 and %d", s.limit)
	}

	return nil
}

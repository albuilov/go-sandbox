package cancellable

import (
	"context"
	"fmt"
)

type Semaphore struct {
	slots chan struct{}
}

func NewSemaphore(limit int) (*Semaphore, error) {
	if limit < 1 {
		return nil, fmt.Errorf("limit must be positive, got %d", limit)
	}

	s := Semaphore{
		slots: make(chan struct{}, limit),
	}

	return &s, nil
}

func (s *Semaphore) Acquire(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.slots <- struct{}{}:
		return nil
	}
}

func (s *Semaphore) TryAcquire() bool {
	select {
	case s.slots <- struct{}{}:
		return true
	default:
		return false
	}
}

func (s *Semaphore) Release() {
	select {
	case <-s.slots:
	default:
		panic("semaphore: release without acquire")
	}
}

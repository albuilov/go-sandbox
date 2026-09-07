package cancellable

import (
	"context"
	"fmt"
)

// Semaphore ограничивает число операций и позволяет отменить ожидание
// Создаем через NewSemaphore. Порядок получения мест не гарантируется
type Semaphore struct {
	slots chan struct{}
}

// NewSemaphore создает семафор с лимитом больше нуля
func NewSemaphore(limit int) (*Semaphore, error) {
	if limit < 1 {
		return nil, fmt.Errorf("limit must be positive, got %d", limit)
	}

	s := Semaphore{
		slots: make(chan struct{}, limit),
	}

	return &s, nil
}

// Acquire ждет свободное место или отмену ctx
// При ошибке место не занято, Release вызывать не нужно
// Если отмена совпала с получением места, метод может вернуть nil
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

// TryAcquire занимает место без ожидания. Если мест нет, возвращает false
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.slots <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release освобождает одно занятое место. Если занятых мест нет, возвращает ошибку
// Освобождаем место ровно один раз после успешного Acquire или TryAcquire
func (s *Semaphore) Release() error {
	select {
	case <-s.slots:
		return nil
	default:
		return fmt.Errorf("semaphore: release without acquire")
	}
}

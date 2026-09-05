package weighted

import (
	"context"
	"fmt"
	"sync"
)

// Semaphore ограничивает общий вес одновременно выполняемых задач.
// Создаем через NewSemaphore и не копируем после начала использования.
// Очереди нет: небольшие задачи могут обгонять большие.
type Semaphore struct {
	mu      sync.Mutex
	limit   int
	used    int
	changed chan struct{}
}

// NewSemaphore создает семафор с общим лимитом больше нуля.
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

// Acquire ждет, пока хватит места для всего веса, или отменят ctx.
// Вес должен быть от 1 до limit. При ошибке разрешения не заняты.
func (s *Semaphore) Acquire(ctx context.Context, weight int) error {
	if err := s.checkWeight(weight); err != nil {
		return fmt.Errorf("acquire error: %w", err)
	}

	for {
		s.mu.Lock()

		if err := ctx.Err(); err != nil {
			s.mu.Unlock()
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

// TryAcquire пытается занять весь вес без ожидания.
func (s *Semaphore) TryAcquire(weight int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkWeight(weight); err != nil {
		return false
	}

	if weight > s.limit-s.used {
		return false
	}

	s.used += weight
	return true
}

// Release возвращает указанный вес. Освобождаем только занятые разрешения.
func (s *Semaphore) Release(weight int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkWeight(weight); err != nil {
		return fmt.Errorf("release error: %w", err)
	}

	if weight > s.used {
		return fmt.Errorf("release weight %d exceeds held weight %d", weight, s.used)
	}

	s.used -= weight

	// Будим всех, кто ждет освобождения разрешений.
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

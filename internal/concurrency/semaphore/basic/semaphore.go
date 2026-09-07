package basic

import "fmt"

// Semaphore ограничивает число одновременно выполняемых операций
// Создаем через NewSemaphore: нулевое значение не готово к работе
type Semaphore struct {
	slots chan struct{}
}

// NewSemaphore создает семафор с лимитом больше нуля
func NewSemaphore(limit int) (*Semaphore, error) {
	if limit < 1 {
		return nil, fmt.Errorf("limit must be positive, got %d", limit)
	}

	s := &Semaphore{
		slots: make(chan struct{}, limit),
	}

	return s, nil
}

// Acquire занимает место. Если все места заняты, ждет
func (s *Semaphore) Acquire() {
	s.slots <- struct{}{}
}

// Release освобождает ранее занятое место
// На каждый Acquire должен приходиться один Release
func (s *Semaphore) Release() {
	<-s.slots
}

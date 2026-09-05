package cancellable

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// RunExample показывает отмену ожидания и завершение уже запущенных задач.
func RunExample() {
	fmt.Println("Semaphore: Cancellable")

	// Короткий таймаут для проверки отмены. Она может не успеть сработать.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	s, err := NewSemaphore(4)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	var wg sync.WaitGroup

	for taskID := range 10 {
		if err := s.Acquire(ctx); err != nil {
			fmt.Println("error:", err)
			break
		}

		wg.Go(func() {
			defer func() {
				if err := s.Release(); err != nil {
					fmt.Println("error:", err)
				}
			}()

			fmt.Printf("task %d: start\n", taskID)
			time.Sleep(time.Second)
			fmt.Printf("task %d: done\n", taskID)
		})
	}

	// Даже при отмене ждем завершения запущенных задач.
	wg.Wait()
}

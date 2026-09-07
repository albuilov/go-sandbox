package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/albuilov/go-sandbox/internal/concurrency/semaphore/cancellable"
)

// RunCancellable показывает отмену ожидания и завершение уже запущенных задач
func RunCancellable(ctx context.Context, limit int) {
	fmt.Println("Semaphore: Cancellable")

	// проверяем отмену коротким таймаутом, но задачи могут успеть завершиться раньше
	ctx, cancel := context.WithTimeout(ctx, 5*time.Millisecond)
	defer cancel()

	s, err := cancellable.NewSemaphore(limit)
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

	// даже при отмене ждем завершения запущенных задач
	wg.Wait()
}

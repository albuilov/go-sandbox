package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/albuilov/go-sandbox/internal/concurrency/semaphore/weighted"
)

// RunWeighted запускает задачи с разным весом при общем лимите
func RunWeighted(ctx context.Context, limit int) {
	fmt.Println("Semaphore: Weighted")

	// проверяем отмену коротким таймаутом, но задачи могут успеть завершиться раньше
	ctx, cancel := context.WithTimeout(ctx, 5*time.Millisecond)
	defer cancel()

	s, err := weighted.NewSemaphore(limit)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	weights := []int{3, 2, 4, 1, 5, 2, 3}

	var wg sync.WaitGroup

	for taskID, weight := range weights {
		// ждем, пока освободится весь нужный вес, и только потом запускаем задачу
		fmt.Printf("task %d: waiting, weight = %d\n", taskID, weight)
		if err := s.Acquire(ctx, weight); err != nil {
			fmt.Println("error:", err)
			break
		}

		wg.Go(func() {
			// возвращаем столько же разрешений, сколько заняли
			defer func() {
				if err := s.Release(weight); err != nil {
					fmt.Println("error:", err)
				}
			}()

			fmt.Printf("task %d: start, weight = %d\n", taskID, weight)
			time.Sleep(time.Second)
			fmt.Printf("task %d: done, weight = %d\n", taskID, weight)
		})
	}

	// даже при отмене ждем завершения запущенных задач
	wg.Wait()
}

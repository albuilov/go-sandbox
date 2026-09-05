package weighted

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// RunExample запускает задачи с разным весом при общем лимите 5.
func RunExample() {
	fmt.Println("Semaphore: Weighted")

	// Короткий таймаут для проверки отмены. Она может не успеть сработать.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	s, err := NewSemaphore(5)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	weights := []int{3, 2, 4, 1, 5, 2, 3}

	var wg sync.WaitGroup

	for taskID, weight := range weights {
		// Ждем, пока освободится весь нужный вес, и только потом запускаем задачу.
		fmt.Printf("task %d: waiting, weight = %d\n", taskID, weight)
		if err := s.Acquire(ctx, weight); err != nil {
			fmt.Println("error:", err)
			break
		}

		wg.Go(func() {
			// Возвращаем столько же разрешений, сколько заняли.
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

	// Даже при отмене ждем завершения запущенных задач.
	wg.Wait()
}

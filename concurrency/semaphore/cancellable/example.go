package cancellable

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func RunExample() {
	fmt.Println("Semaphore: cancellable")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	s, err := NewSemaphore(4)
	if err != nil {
		fmt.Printf("Semaphore cancellable error: %v\n", err)
		return
	}

	var wg sync.WaitGroup

	for taskID := range 10 {
		if err := s.Acquire(ctx); err != nil {
			fmt.Printf("Semaphore Acquire error: %v\n", err)
			break
		}

		wg.Go(func() {
			defer s.Release()

			fmt.Printf("task %d: start\n", taskID)
			time.Sleep(time.Second)
			fmt.Printf("task %d: done\n", taskID)
		})
	}

	// Даже при отмене ждем завершения запущенных задач.
	wg.Wait()
}

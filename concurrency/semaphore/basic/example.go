package basic

import (
	"fmt"
	"sync"
	"time"
)

func RunExample() {
	fmt.Println("Semaphore: Basic")

	s, err := NewSemaphore(4)
	if err != nil {
		fmt.Printf("Semaphore basic error: %v\n", err)
		return
	}

	var wg sync.WaitGroup

	for taskID := range 10 {
		s.Acquire()

		wg.Go(func() {
			defer s.Release()

			fmt.Printf("task %d: start\n", taskID)
			time.Sleep(time.Second)
			fmt.Printf("task %d: done\n", taskID)
		})
	}

	// Ждем завершения запущенных задач.
	wg.Wait()
}

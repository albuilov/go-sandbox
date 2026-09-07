package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/albuilov/go-sandbox/internal/concurrency/semaphore/basic"
)

// RunBasic запускает десять задач с заданным лимитом одновременных операций
func RunBasic(limit int) {
	fmt.Println("Semaphore: Basic")

	s, err := basic.NewSemaphore(limit)
	if err != nil {
		fmt.Println("error:", err)
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

	// ждем завершения запущенных задач
	wg.Wait()
}

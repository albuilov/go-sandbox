package results

import (
	"fmt"
	"sync"
)

// WorkerPool считает и выводит сумму квадратов.
// workerCount должен быть больше нуля.
func WorkerPool(workerCount int, nums []int) (<-chan int, error) {
	if workerCount < 1 {
		return nil, fmt.Errorf("workerCount must be positive, got %d", workerCount)
	}

	out := make(chan int)
	jobs := make(chan int)

	var wg sync.WaitGroup

	for range workerCount {
		wg.Go(func() {
			for num := range jobs {
				out <- num * num
			}
		})
	}

	go func() {
		defer close(out)
		wg.Wait()
	}()

	// Отправляем задания параллельно с чтением результатов.
	go func() {
		defer close(jobs)
		for _, num := range nums {
			jobs <- num
		}
	}()

	return out, nil
}

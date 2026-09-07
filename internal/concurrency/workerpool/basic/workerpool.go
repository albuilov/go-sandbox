package basic

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// WorkerPool считает и возвращает сумму квадратов
// workerCount должен быть больше нуля
func WorkerPool(workerCount int, nums []int) (int, error) {
	if workerCount < 1 {
		return 0, fmt.Errorf("workerCount must be positive, got %d", workerCount)
	}

	jobs := make(chan int)
	// несколько воркеров обновляют общую сумму,
	// поэтому используем атомарное сложение
	sum := atomic.Int64{}

	var wg sync.WaitGroup

	for range workerCount {
		wg.Go(func() {
			for num := range jobs {
				sum.Add(int64(num * num))
			}
		})
	}

	for _, num := range nums {
		jobs <- num
	}

	close(jobs)
	wg.Wait()

	return int(sum.Load()), nil
}

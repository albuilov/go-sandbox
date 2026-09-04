package basic

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Run демонстрирует суммирование квадратов общим атомарным счётчиком.
func Run(workerCount int, nums []int) error {
	fmt.Println("Worker Pool basic")

	res, err := pool(workerCount, nums)
	if err != nil {
		return fmt.Errorf("run basic worker pool: %w", err)
	}

	fmt.Println("Total sum: ", res)

	return nil
}

// pool возвращает сумму квадратов после завершения всех воркеров.
// workerCount должен быть положительным.
func pool(workerCount int, nums []int) (int, error) {
	if workerCount < 1 {
		return 0, fmt.Errorf("workerCount must be positive, got %d", workerCount)
	}

	jobs := make(chan int)
	// Несколько воркеров обновляют общую сумму, поэтому используем атомарное сложение.
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

package results

import (
	"fmt"
	"sync"
)

// Run собирает квадраты из канала и выводит их сумму.
func Run(workerCount int, nums []int) error {
	fmt.Println("Worker Pool: Results")

	res, err := pool(workerCount, nums)
	if err != nil {
		return fmt.Errorf("run result channel worker pool: %w", err)
	}

	var sum int
	for num := range res {
		sum += num
	}

	fmt.Println("Total sum:", sum)

	return nil
}

// pool запускает workerCount воркеров и возвращает канал квадратов переданных чисел.
// Порядок результатов не гарантируется.
// Нужно прочитать канал до конца, иначе воркеры останутся ждать отправки.
// workerCount должен быть положительным.
func pool(workerCount int, nums []int) (<-chan int, error) {
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

	// Отправляем задания в отдельной горутине, чтобы сразу вернуть канал результатов.
	go func() {
		defer close(jobs)
		for _, num := range nums {
			jobs <- num
		}
	}()

	return out, nil
}

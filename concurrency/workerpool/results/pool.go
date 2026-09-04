package results

import (
	"fmt"
	"sync"
)

// Run демонстрирует сбор квадратов через канал и суммирование получателем.
func Run(workerCount int, nums []int) error {
	fmt.Println("Worker Pool with channel result")

	res, err := pool(workerCount, nums)
	if err != nil {
		return fmt.Errorf("run result channel worker pool: %w", err)
	}

	var sum int
	for num := range res {
		sum += num
	}

	fmt.Println("Total sum: ", sum)

	return nil
}

// pool запускает workerCount воркеров и возвращает канал квадратов переданных чисел.
// Порядок результатов не гарантируется.
// Вызывающий должен прочитать канал до конца, иначе воркеры заблокируются.
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

	// Отправляем задания асинхронно, чтобы вызывающий мог начать читать результаты.
	go func() {
		defer close(jobs)
		for _, num := range nums {
			jobs <- num
		}
	}()

	return out, nil
}

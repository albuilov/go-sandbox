package cancellable

import (
	"context"
	"fmt"
	"sync"
)

// WorkerPool считает и выводит сумму квадратов.
// workerCount должен быть больше нуля.
// При отмене возвращает ошибку вместо вывода полной суммы.
func WorkerPool(ctx context.Context, workerCount int, nums []int) error {
	if workerCount < 1 {
		return fmt.Errorf("workerCount must be positive, got %d", workerCount)
	}

	out := make(chan int)
	jobs := make(chan int)

	var wg sync.WaitGroup

	for range workerCount {
		wg.Go(func() {
			for num := range jobs {
				select {
				case <-ctx.Done():
					return
				case out <- num * num:
				}
			}
		})
	}

	go func() {
		defer close(out)
		wg.Wait()
	}()

	go func() {
		// Закрываем jobs и при отмене,
		// чтобы воркеры не ждали новых заданий.
		defer close(jobs)

		for _, num := range nums {
			select {
			case <-ctx.Done():
				return
			case jobs <- num:
			}
		}
	}()

	var sum int
	for num := range out {
		sum += num
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("run cancellable worker pool: %w", err)
	}

	fmt.Println("Total sum:", sum)
	return nil
}

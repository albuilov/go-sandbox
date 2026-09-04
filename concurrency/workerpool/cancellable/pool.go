package cancellable

import (
	"context"
	"fmt"
	"sync"
)

// Run демонстрирует пул с отменой через контекст.
// При отмене возвращает ошибку вместо вывода полной суммы.
func Run(ctx context.Context, workerCount int, nums []int) error {
	fmt.Println("Worker Pool with context and channel result")

	res, err := pool(ctx, workerCount, nums)
	if err != nil {
		return fmt.Errorf("run cancellable worker pool: %w", err)
	}

	var sum int
	for num := range res {
		sum += num
	}

	if err = ctx.Err(); err != nil {
		return fmt.Errorf("run cancellable worker pool: %w", err)
	}

	fmt.Println("Total sum: ", sum)

	return nil
}

// pool запускает workerCount воркеров и возвращает канал квадратов чисел.
// workerCount должен быть положительным. Порядок результатов не гарантируется.
// Вызывающий должен читать канал до закрытия или отменить ctx.
// При отмене часть результатов может не поступить.
// Канал результатов закрывается после завершения всех воркеров.
func pool(ctx context.Context, workerCount int, nums []int) (<-chan int, error) {
	if workerCount < 1 {
		return nil, fmt.Errorf("workerCount must be positive, got %d", workerCount)
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
		// Закрываем jobs и при отмене, освобождая воркеров, ожидающих задания.
		defer close(jobs)

		for _, num := range nums {
			select {
			case <-ctx.Done():
				return
			case jobs <- num:
			}
		}
	}()

	return out, nil
}

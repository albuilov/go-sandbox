package cancellable

import (
	"context"
	"fmt"
	"sync"
)

// WorkerPool запускает воркеров и возвращает канал с квадратами чисел
// workerCount должен быть больше нуля. Порядок результатов может меняться
// Читаем канал до конца или отменяем ctx, чтобы воркеры не ждали отправки
// После возврата функции отмену проверяем через ctx.Err()
func WorkerPool(ctx context.Context, workerCount int, nums []int) (<-chan int, error) {
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

	// закрываем результаты только после завершения всех воркеров
	go func() {
		defer close(out)
		wg.Wait()
	}()

	go func() {
		// закрываем jobs и при отмене,
		// чтобы воркеры не ждали новых заданий
		defer close(jobs)

		for _, num := range nums {
			select {
			case <-ctx.Done():
				return
			case jobs <- num:
			}
		}
	}()

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("run cancellable worker pool: %w", err)
	}

	return out, nil
}

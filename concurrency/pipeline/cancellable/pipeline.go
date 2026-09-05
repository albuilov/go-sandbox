package cancellable

import (
	"context"
	"fmt"
	"sync"
)

// Pipeline запускает цепочку с отменой: числа → квадраты → четные числа.
// Все этапы используют один контекст, чтобы при отмене остановить всю цепочку.
func Pipeline(ctx context.Context, workerCount int, nums []int) error {
	// Проверяем число воркеров до запуска генератора.
	if workerCount < 1 {
		return fmt.Errorf("workerCount must be positive, got %d", workerCount)
	}

	numbers := generate(ctx, nums)
	squares := square(ctx, workerCount, numbers)
	evens := filterEven(ctx, squares)

	for num := range evens {
		fmt.Println("->", num)
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("run cancellable pipeline: %w", err)
	}

	return nil
}

// generate отправляет числа в канал, пока они не закончатся или не отменят ctx.
// Закрываем канал и при отмене, чтобы следующий этап не ждал новых чисел.
func generate(ctx context.Context, nums []int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, num := range nums {
			select {
			case <-ctx.Done():
				return
			case out <- num:
			}
		}
	}()

	return out
}

// square считает квадраты в нескольких воркерах.
// workerCount должен быть больше нуля. Результаты могут прийти в другом порядке.
// При отмене можно выйти, даже если следующий этап перестал читать результаты.
func square(ctx context.Context, workerCount int, in <-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup

	for range workerCount {
		wg.Go(func() {
			for num := range in {
				select {
				case <-ctx.Done():
					return
				case out <- num * num:
				}
			}
		})
	}

	// Закрываем канал, когда все воркеры закончили отправлять результаты.
	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}

// filterEven пропускает только четные числа в том порядке, в котором получил их.
// При отмене перестаем отправлять результаты.
func filterEven(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for num := range in {
			if num%2 != 0 {
				continue
			}

			select {
			case <-ctx.Done():
				return
			case out <- num:
			}
		}
	}()

	return out
}

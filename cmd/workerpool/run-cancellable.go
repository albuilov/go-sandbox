package main

import (
	"context"
	"fmt"
	"time"

	"github.com/albuilov/go-sandbox/internal/concurrency/workerpool/cancellable"
)

// RunCancellable демонстрирует пул с отменой через контекст
// При отмене выводит ошибку вместо полной суммы
func RunCancellable(ctx context.Context, workerCount int, nums []int) {
	fmt.Println("Worker Pool: Cancellable")

	// проверяем отмену коротким таймаутом, но задачи могут успеть завершиться раньше
	ctx, cancel := context.WithTimeout(ctx, 5*time.Microsecond)
	defer cancel()

	val, err := cancellable.WorkerPool(ctx, workerCount, nums)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	var sum int
	for num := range val {
		sum += num
	}

	// после закрытия канала проверяем, не было ли отмены
	if err := ctx.Err(); err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("Total sum:", sum)
}

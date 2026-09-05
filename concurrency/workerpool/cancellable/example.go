package cancellable

import (
	"context"
	"fmt"
	"time"
)

// RunExample демонстрирует пул с отменой через контекст.
// При отмене выводит ошибку вместо полной суммы.
func RunExample(workerCount int, nums []int) {
	fmt.Println("Worker Pool: Cancellable")

	// Короткий таймаут для проверки отмены. Она может не успеть сработать.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Microsecond)
	defer cancel()

	if err := WorkerPool(ctx, workerCount, nums); err != nil {
		fmt.Println("error:", err)
	}
}

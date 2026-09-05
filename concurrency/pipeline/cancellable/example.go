package cancellable

import (
	"context"
	"fmt"
	"time"
)

// RunExample запускает pipeline с коротким таймаутом.
func RunExample(workerCount int, nums []int) {
	fmt.Println("Pipeline: Cancellable")

	// Короткий таймаут для проверки отмены. Она может не успеть сработать.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Microsecond)
	defer cancel()

	if err := Pipeline(ctx, workerCount, nums); err != nil {
		fmt.Println("error:", err)
		return
	}
}

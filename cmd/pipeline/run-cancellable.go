package main

import (
	"context"
	"fmt"
	"time"

	"github.com/albuilov/go-sandbox/internal/concurrency/pipeline/cancellable"
)

// RunCancellable запускает pipeline с коротким таймаутом
func RunCancellable(ctx context.Context, workerCount int, nums []int) {
	fmt.Println("Pipeline: Cancellable")

	// проверяем отмену коротким таймаутом, но задачи могут успеть завершиться раньше
	ctx, cancel := context.WithTimeout(ctx, 5*time.Microsecond)
	defer cancel()

	evens, err := cancellable.NewPipeline(ctx, workerCount, nums)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	var sum int
	for num := range evens {
		sum += num
	}

	// после закрытия канала проверяем, не было ли отмены
	if err := ctx.Err(); err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("Total sum:", sum)
}

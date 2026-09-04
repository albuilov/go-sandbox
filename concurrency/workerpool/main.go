package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/albuilov/go-sandbox/concurrency/workerpool/basic"
	"github.com/albuilov/go-sandbox/concurrency/workerpool/cancellable"
	"github.com/albuilov/go-sandbox/concurrency/workerpool/results"
)

func main() {
	var err error

	workerCount := 3
	nums := generateNumbers(10)

	if err = basic.Run(workerCount, nums); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()

	if err = results.Run(workerCount, nums); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()

	// Короткий таймаут для демонстрации отмены; результат зависит от планирования.
	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()

	if err = cancellable.Run(ctx, workerCount, nums); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func generateNumbers(count int) []int {
	if count < 0 {
		count = 0
	}

	nums := make([]int, count)
	for i := range count {
		nums[i] = randomNumber(10, 99)
	}

	return nums
}

func randomNumber(min, max int) int {
	return rand.Intn(max-min+1) + min
}

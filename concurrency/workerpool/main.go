package main

import (
	"context"
	"fmt"
	"os"

	"github.com/albuilov/go-sandbox/concurrency/workerpool/basic"
	"github.com/albuilov/go-sandbox/concurrency/workerpool/cancellable"
	"github.com/albuilov/go-sandbox/concurrency/workerpool/results"
	"github.com/albuilov/go-sandbox/internal/number"
)

// Раздаём числа воркерам, считаем квадраты и выводим их сумму.
func main() {
	workerCount := 3
	nums := number.GenerateNumbers(10, 10, 99)

	if err := basic.Run(workerCount, nums); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "worker pool error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()

	if err := results.Run(workerCount, nums); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "worker pool error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()

	// Пока читаем все результаты. Для проверки отмены нужно вызвать cancel раньше.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := cancellable.Run(ctx, workerCount, nums); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "worker pool error: %v\n", err)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/albuilov/go-sandbox/concurrency/pipeline/basic"
	"github.com/albuilov/go-sandbox/concurrency/pipeline/cancellable"
	"github.com/albuilov/go-sandbox/internal/number"
)

// Отправляем случайные числа от 10 до 99 в канал.
// Считаем квадраты, оставляем чётные и выводим результат.
func main() {
	workerCount := 3
	nums := number.GenerateNumbers(10, 10, 99)

	if err := basic.Run(workerCount, nums); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "pipeline error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()

	// Пока читаем все результаты. Для проверки отмены нужно вызвать cancel раньше.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := cancellable.Run(ctx, workerCount, nums); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "pipeline error: %v\n", err)
		os.Exit(1)
	}
}

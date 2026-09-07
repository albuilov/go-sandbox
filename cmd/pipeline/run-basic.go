package main

import (
	"fmt"

	"github.com/albuilov/go-sandbox/internal/concurrency/pipeline/basic"
)

// RunBasic демонстрирует работу pipeline:
// числа → квадраты → четные числа
func RunBasic(workerCount int, nums []int) {
	fmt.Println("Pipeline: Basic")

	evens, err := basic.NewPipeline(workerCount, nums)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	var sum int
	for num := range evens {
		sum += num
	}

	fmt.Println("Total sum:", sum)
}

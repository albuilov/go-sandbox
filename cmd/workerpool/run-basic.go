package main

import (
	"fmt"

	"github.com/albuilov/go-sandbox/internal/concurrency/workerpool/basic"
)

// RunBasic считает сумму квадратов с помощью общего атомарного счетчика
func RunBasic(workerCount int, nums []int) {
	fmt.Println("Worker Pool: Basic")

	val, err := basic.WorkerPool(workerCount, nums)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("Total sum:", val)
}

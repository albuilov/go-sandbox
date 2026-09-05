package basic

import (
	"fmt"
)

// RunExample считает сумму квадратов с помощью общего атомарного счетчика.
func RunExample(workerCount int, nums []int) {
	fmt.Println("Worker Pool: Basic")

	if err := WorkerPool(workerCount, nums); err != nil {
		fmt.Println("error:", err)
	}
}

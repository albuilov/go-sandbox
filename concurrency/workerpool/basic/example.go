package basic

import (
	"fmt"
)

// RunExample считает сумму квадратов с помощью общего атомарного счетчика.
func RunExample(workerCount int, nums []int) {
	fmt.Println("Worker Pool: Basic")

	val, err := WorkerPool(workerCount, nums)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("Total sum:", val)
}

package results

import (
	"fmt"
)

// RunExample собирает квадраты из канала и выводит их сумму.
func RunExample(workerCount int, nums []int) {
	fmt.Println("Worker Pool: Results")

	if err := WorkerPool(workerCount, nums); err != nil {
		fmt.Println("error:", err)
	}
}

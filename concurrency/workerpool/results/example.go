package results

import (
	"fmt"
)

// RunExample собирает квадраты из канала и выводит их сумму.
func RunExample(workerCount int, nums []int) {
	fmt.Println("Worker Pool: Results")

	val, err := WorkerPool(workerCount, nums)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	var sum int
	for num := range val {
		sum += num
	}

	fmt.Println("Total sum:", sum)
}

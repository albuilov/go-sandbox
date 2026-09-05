package basic

import (
	"fmt"
)

// RunExample демонстрирует работу pipeline:
// числа → квадраты → четные числа.
func RunExample(workerCount int, nums []int) {
	fmt.Println("Pipeline: Basic")

	if err := Pipeline(workerCount, nums); err != nil {
		fmt.Println("error:", err)
		return
	}
}

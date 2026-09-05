package number

import "math/rand"

// GenerateNumbers создаёт count случайных чисел от min до max включительно.
// Если count меньше или равен нулю, возвращает пустой слайс.
func GenerateNumbers(count, min, max int) []int {
	if count < 0 {
		count = 0
	}

	nums := make([]int, count)
	for i := range count {
		nums[i] = RandomNumber(min, max)
	}

	return nums
}

// RandomNumber возвращает случайное число от min до max включительно.
func RandomNumber(min, max int) int {
	return rand.Intn(max-min+1) + min
}

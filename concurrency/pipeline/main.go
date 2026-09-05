package main

import (
	"github.com/albuilov/go-sandbox/concurrency/pipeline/basic"
	"github.com/albuilov/go-sandbox/concurrency/pipeline/cancellable"
	"github.com/albuilov/go-sandbox/internal/number"
)

// Отправляем случайные числа от 10 до 99 в канал.
// Считаем квадраты, оставляем четные и выводим результат.
func main() {
	workerCount := 3
	nums := number.GenerateNumbers(10, 10, 99)

	basic.RunExample(workerCount, nums)
	cancellable.RunExample(workerCount, nums)
}

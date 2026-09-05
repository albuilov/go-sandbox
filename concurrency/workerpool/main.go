package main

import (
	"github.com/albuilov/go-sandbox/concurrency/workerpool/basic"
	"github.com/albuilov/go-sandbox/concurrency/workerpool/cancellable"
	"github.com/albuilov/go-sandbox/concurrency/workerpool/results"
	"github.com/albuilov/go-sandbox/internal/number"
)

// Раздаем числа воркерам, считаем квадраты и выводим их сумму.
func main() {
	workerCount := 3
	nums := number.GenerateNumbers(10, 10, 99)

	basic.RunExample(workerCount, nums)
	results.RunExample(workerCount, nums)
	cancellable.RunExample(workerCount, nums)
}

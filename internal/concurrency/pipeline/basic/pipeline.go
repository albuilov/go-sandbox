package basic

import (
	"fmt"
	"sync"
)

// NewPipeline запускает цепочку: числа → квадраты → четные числа
// Читаем результат до конца, чтобы не оставить горутины ждать отправки
func NewPipeline(workerCount int, nums []int) (<-chan int, error) {
	// проверяем число воркеров до запуска генератора
	if workerCount < 1 {
		return nil, fmt.Errorf("workerCount must be positive, got %d", workerCount)
	}

	numbers := generate(nums)
	squares := square(workerCount, numbers)
	evens := filterEven(squares)

	return evens, nil
}

// generate отправляет числа из слайса в канал, затем закрывает его
func generate(nums []int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, num := range nums {
			out <- num
		}
	}()

	return out
}

// square считает квадраты в нескольких воркерах
// workerCount должен быть больше нуля. Результаты могут прийти в другом порядке
func square(workerCount int, in <-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup

	for range workerCount {
		wg.Go(func() {
			for num := range in {
				out <- num * num
			}
		})
	}

	// закрываем канал, когда все воркеры закончили отправлять результаты
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// filterEven пропускает только четные числа в том порядке, в котором получил их
func filterEven(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for num := range in {
			if num%2 != 0 {
				continue
			}
			out <- num
		}
	}()

	return out
}

package main

import (
	"context"
	"fmt"

	"github.com/albuilov/go-sandbox/pkg/number"
)

// Отправляем случайные числа от 10 до 99 в канал
// Считаем квадраты, оставляем четные и выводим результат
func main() {
	wCtn := 3
	nums := number.GenerateNumbers(10, 10, 99)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	RunBasic(wCtn, nums)
	RunCancellable(ctx, wCtn, nums)

	fmt.Println("The End!")
}

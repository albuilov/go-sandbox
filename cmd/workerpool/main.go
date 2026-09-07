package main

import (
	"context"
	"fmt"

	"github.com/albuilov/go-sandbox/pkg/number"
)

func main() {
	wCtn := 3
	nums := number.GenerateNumbers(10, 10, 99)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	RunBasic(wCtn, nums)
	RunCancellable(ctx, wCtn, nums)

	fmt.Println("The End!")
}

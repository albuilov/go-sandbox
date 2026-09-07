package main

import (
	"context"
	"fmt"
)

func main() {
	limit := 4

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	RunBasic(limit)
	RunCancellable(ctx, limit)
	RunWeighted(ctx, limit)

	fmt.Println("The End!")
}

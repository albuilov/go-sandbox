package main

import (
	"fmt"

	"github.com/albuilov/go-sandbox/concurrency/semaphore/basic"
	"github.com/albuilov/go-sandbox/concurrency/semaphore/cancellable"
	"github.com/albuilov/go-sandbox/concurrency/semaphore/weighted"
)

func main() {
	basic.RunExample()
	fmt.Println()
	cancellable.RunExample()
	fmt.Println()
	weighted.RunExample()
}

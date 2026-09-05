package main

import (
	"github.com/albuilov/go-sandbox/concurrency/semaphore/basic"
	"github.com/albuilov/go-sandbox/concurrency/semaphore/cancellable"
	"github.com/albuilov/go-sandbox/concurrency/semaphore/weighted"
)

func main() {
	basic.RunExample()
	cancellable.RunExample()
	weighted.RunExample()
}

package basic

import (
	"fmt"
	"testing"
	"testing/synctest"
)

func TestWorkerPool(t *testing.T) {
	tests := []struct {
		name    string
		workers int
		nums    []int
		want    int
	}{
		{"nil input", 3, nil, 0},
		{"empty input", 3, []int{}, 0},
		{"one worker", 1, []int{1, 2, 3, 4}, 30},
		{"several workers", 3, []int{1, 2, 3, 4, 5}, 55},
		{"more workers than jobs", 8, []int{2, 3}, 13},
		{"negative and repeated numbers", 3, []int{-3, -2, 0, 2, 2}, 21},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				got, err := WorkerPool(tc.workers, tc.nums)
				if err != nil {
					t.Fatal(err)
				}
				if got != tc.want {
					t.Fatalf("sum = %d; want %d", got, tc.want)
				}
			})
		})
	}
}

func TestWorkerPoolInvalidWorkerCount(t *testing.T) {
	for _, workers := range []int{-1, 0} {
		t.Run(fmt.Sprint(workers), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				if _, err := WorkerPool(workers, []int{1, 2}); err == nil {
					t.Fatal("expected invalid worker count error")
				}
			})
		})
	}
}

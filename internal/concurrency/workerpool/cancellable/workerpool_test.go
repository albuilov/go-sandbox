package cancellable

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"testing/synctest"
	"time"
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
				out, err := WorkerPool(context.Background(), tc.workers, tc.nums)
				if err != nil {
					t.Fatal(err)
				}
				var got int
				for num := range out {
					got += num
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
				if _, err := WorkerPool(context.Background(), workers, []int{1, 2}); err == nil {
					t.Fatal("expected invalid worker count error")
				}
			})
		})
	}
}

func TestWorkerPoolCancelledContext(t *testing.T) {
	for _, scenario := range []string{"cancel", "deadline"} {
		t.Run(scenario, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				want := context.Canceled
				if scenario == "cancel" {
					cancel()
				} else {
					// в synctest время виртуальное
					time.Sleep(time.Second)
					want = context.DeadlineExceeded
				}
				out, err := WorkerPool(ctx, 3, []int{2, 4, 6})
				if !errors.Is(err, want) {
					t.Fatalf("error = %v; want %v", err, want)
				}
				if out != nil {
					t.Fatal("expected nil channel for cancelled context")
				}
			})
		})
	}
}

func TestCancelWithoutReadingResults(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		out, err := WorkerPool(ctx, 3, []int{2, 4, 6, 8, 10, 12})
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := <-out; !ok {
			t.Fatal("worker pool closed before first result")
		}

		// больше не читаем. Даем воркерам заблокироваться на отправке
		synctest.Wait()
		cancel()
		synctest.Wait()

		// все горутины пула должны закончить работу без помощи получателя
		select {
		case _, ok := <-out:
			if ok {
				t.Fatal("output still open after cancellation")
			}
		default:
			t.Fatal("output did not close after cancellation")
		}
	})
}

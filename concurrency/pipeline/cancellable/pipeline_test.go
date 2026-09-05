package cancellable

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
	"testing/synctest"
	"time"
)

func TestPipeline(t *testing.T) {
	tests := []struct {
		name    string
		workers int
		nums    []int
		want    []int
	}{
		{"nil input", 3, nil, nil},
		{"empty input", 3, []int{}, nil},
		{"one worker", 1, []int{1, 2, 3, 4, 5, 6}, []int{4, 16, 36}},
		{"several workers", 3, []int{1, 2, 3, 4, 5, 6}, []int{4, 16, 36}},
		{"more workers than jobs", 8, []int{2}, []int{4}},
		{"only odd numbers", 3, []int{-3, 1, 5}, nil},
		{"negative and repeated numbers", 3, []int{-4, -2, 0, 2, 2}, []int{0, 4, 4, 4, 16}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output := captureOutput(t, func() {
				synctest.Test(t, func(t *testing.T) {
					if err := Pipeline(context.Background(), tc.workers, tc.nums); err != nil {
						t.Fatal(err)
					}
				})
			})

			var got []int
			for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
				if line == "" {
					continue
				}

				if !strings.HasPrefix(line, "-> ") {
					t.Fatalf("unexpected output: %q", line)
				}

				num, err := strconv.Atoi(strings.TrimPrefix(line, "-> "))
				if err != nil {
					t.Fatal(err)
				}

				got = append(got, num)
			}

			// Воркеры могут поменять порядок, но не значения и их количество.
			slices.Sort(got)
			if !slices.Equal(got, tc.want) {
				t.Fatalf("results = %v; want %v", got, tc.want)
			}
		})
	}
}

func TestPipelineInvalidWorkerCount(t *testing.T) {
	for _, workers := range []int{-1, 0} {
		t.Run(fmt.Sprint(workers), func(t *testing.T) {
			output := captureOutput(t, func() {
				synctest.Test(t, func(t *testing.T) {
					if err := Pipeline(context.Background(), workers, []int{1, 2}); err == nil {
						t.Fatal("expected invalid worker count error")
					}
				})
			})

			if output != "" {
				t.Fatalf("invalid input produced output: %q", output)
			}
		})
	}
}

func TestGeneratePreservesOrder(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		nums := []int{3, -2, 0, 3}
		var got []int
		for num := range generate(context.Background(), nums) {
			got = append(got, num)
		}

		if !slices.Equal(got, nums) {
			t.Fatalf("generate = %v; want %v", got, nums)
		}
	})
}

func TestFilterEvenPreservesOrder(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		in := make(chan int)

		go func() {
			defer close(in)
			for _, num := range []int{8, 3, -2, 8, 0, 5} {
				in <- num
			}
		}()

		var got []int
		for num := range filterEven(context.Background(), in) {
			got = append(got, num)
		}

		want := []int{8, -2, 8, 0}
		if !slices.Equal(got, want) {
			t.Fatalf("filterEven = %v; want %v", got, want)
		}
	})
}

func TestPipelineCancelledContext(t *testing.T) {
	for _, scenario := range []string{"cancel", "deadline"} {
		t.Run(scenario, func(t *testing.T) {
			captureOutput(t, func() {
				synctest.Test(t, func(t *testing.T) {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					want := context.Canceled
					if scenario == "cancel" {
						cancel()
					} else {
						time.Sleep(time.Second)
						want = context.DeadlineExceeded
					}

					if err := Pipeline(ctx, 3, []int{2, 4, 6}); !errors.Is(err, want) {
						t.Fatalf("Pipeline = %v; want %v", err, want)
					}
				})
			})
		})
	}
}

func TestCancelWithoutReadingResults(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		out := filterEven(ctx, square(ctx, 3, generate(ctx, []int{2, 4, 6, 8, 10, 12})))
		if _, ok := <-out; !ok {
			t.Fatal("pipeline closed before first result")
		}

		// Больше не читаем. Даем этапам заблокироваться на отправке.
		synctest.Wait()
		cancel()
		synctest.Wait()

		// Все этапы должны закончить работу без помощи получателя.
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

// Перехватываем вывод реализации. Эти тесты нельзя запускать через t.Parallel.
func captureOutput(t *testing.T, run func()) string {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), "output-*")
	if err != nil {
		t.Fatal(err)
	}

	original := os.Stdout
	os.Stdout = file
	defer func() {
		os.Stdout = original
		file.Close()
	}()

	run()
	if _, err := file.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	output, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	return string(output)
}

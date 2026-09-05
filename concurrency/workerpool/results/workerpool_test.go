package results

import (
	"fmt"
	"io"
	"os"
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
			output := captureOutput(t, func() {
				// Test дождется всех горутин и обнаружит зависание внутри цепочки.
				synctest.Test(t, func(t *testing.T) {
					if err := WorkerPool(tc.workers, tc.nums); err != nil {
						t.Fatal(err)
					}
				})
			})

			if want := fmt.Sprintf("Total sum: %d\n", tc.want); output != want {
				t.Fatalf("output = %q; want %q", output, want)
			}
		})
	}
}

func TestWorkerPoolInvalidWorkerCount(t *testing.T) {
	for _, workers := range []int{-1, 0} {
		t.Run(fmt.Sprint(workers), func(t *testing.T) {
			output := captureOutput(t, func() {
				synctest.Test(t, func(t *testing.T) {
					if err := WorkerPool(workers, []int{1, 2}); err == nil {
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

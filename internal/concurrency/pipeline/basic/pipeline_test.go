package basic

import (
	"fmt"
	"slices"
	"testing"
	"testing/synctest"
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
			synctest.Test(t, func(t *testing.T) {
				out, err := NewPipeline(tc.workers, tc.nums)
				if err != nil {
					t.Fatal(err)
				}
				var got []int
				for num := range out {
					got = append(got, num)
				}
				// воркеры могут поменять порядок, но не значения и их количество
				slices.Sort(got)
				if !slices.Equal(got, tc.want) {
					t.Fatalf("results = %v; want %v", got, tc.want)
				}
			})
		})
	}
}

func TestPipelineInvalidWorkerCount(t *testing.T) {
	for _, workers := range []int{-1, 0} {
		t.Run(fmt.Sprint(workers), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				if _, err := NewPipeline(workers, []int{1, 2}); err == nil {
					t.Fatal("expected invalid worker count error")
				}
			})
		})
	}
}

func TestGeneratePreservesOrder(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		nums := []int{3, -2, 0, 3}
		var got []int
		for num := range generate(nums) {
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
		for num := range filterEven(in) {
			got = append(got, num)
		}
		want := []int{8, -2, 8, 0}
		if !slices.Equal(got, want) {
			t.Fatalf("filterEven = %v; want %v", got, want)
		}
	})
}

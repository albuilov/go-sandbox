package lru

import (
	"fmt"
	"sync"
	"testing"
)

func newTestCache[K comparable, V any](t *testing.T, capacity int) Cache[K, V] {
	t.Helper()
	c, err := NewCache[K, V](capacity)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func checkGet[K comparable, V comparable](t *testing.T, c Cache[K, V], key K, want V, exists bool) {
	t.Helper()
	if got, ok := c.Get(key); got != want || ok != exists {
		t.Fatalf("Get(%v) = (%v, %v); want (%v, %v)", key, got, ok, want, exists)
	}
}

func TestNewCache(t *testing.T) {
	for _, capacity := range []int{-1, 0, 1, 3} {
		t.Run(fmt.Sprint(capacity), func(t *testing.T) {
			c, err := NewCache[int, string](capacity)
			if capacity < 1 {
				if err == nil || c != nil {
					t.Fatalf("NewCache(%d) = (%v, %v); want nil and error", capacity, c, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if c.Len() != 0 {
				t.Fatalf("Len() = %d; want 0", c.Len())
			}
			checkGet(t, c, 0, "", false)
		})
	}
}

// проверяем, что чтение и обновление меняют порядок вытеснения
func TestEviction(t *testing.T) {
	for _, tc := range []struct {
		name  string
		touch func(*testing.T, Cache[string, int])
		lost  string
	}{
		{"oldest", func(t *testing.T, c Cache[string, int]) {}, "A"},
		{"get_tail", func(t *testing.T, c Cache[string, int]) { checkGet(t, c, "A", 1, true) }, "B"},
		{"get_head", func(t *testing.T, c Cache[string, int]) { checkGet(t, c, "B", 2, true) }, "A"},
		{"get_missing", func(t *testing.T, c Cache[string, int]) { checkGet(t, c, "X", 0, false) }, "A"},
		{"update_tail", func(t *testing.T, c Cache[string, int]) { c.Set("A", 10) }, "B"},
		{"update_head", func(t *testing.T, c Cache[string, int]) { c.Set("B", 20) }, "A"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := newTestCache[string, int](t, 2)
			c.Set("A", 1)
			c.Set("B", 2)
			tc.touch(t, c)
			if c.Len() != 2 {
				t.Fatalf("Len() = %d; want 2", c.Len())
			}
			c.Set("C", 3)
			checkGet(t, c, tc.lost, 0, false)
			checkGet(t, c, "C", 3, true)
			key, value := "A", 1
			if tc.lost == "A" {
				key, value = "B", 2
			}
			if tc.name == "update_tail" || tc.name == "update_head" {
				value *= 10
			}
			checkGet(t, c, key, value, true)
			if c.Len() != 2 {
				t.Fatalf("Len() = %d; want 2", c.Len())
			}
		})
	}
}

func TestSingleCapacity(t *testing.T) {
	c := newTestCache[int, string](t, 1)
	for key := range 5 {
		c.Set(key, "old")
		c.Set(key, "new")
		checkGet(t, c, key, "new", true)
		checkGet(t, c, key-1, "", false)
		if c.Len() != 1 {
			t.Fatalf("Len() = %d; want 1", c.Len())
		}
	}
}

func TestRemove(t *testing.T) {
	for _, tc := range []struct {
		name string
		size int
		key  int
	}{
		{"empty", 0, 0}, {"missing", 3, 9}, {"only", 1, 0},
		{"head", 3, 2}, {"middle", 3, 1}, {"tail", 3, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := newTestCache[int, int](t, 3)
			for key := range tc.size {
				c.Set(key, key+10)
			}
			c.Remove(tc.key)
			c.Remove(tc.key) // повторное удаление ничего не меняет
			checkGet(t, c, tc.key, 0, false)
			want := tc.size
			if tc.key < tc.size {
				want--
			}
			if c.Len() != want {
				t.Fatalf("Len() = %d; want %d", c.Len(), want)
			}
			for key := range tc.size {
				if key != tc.key {
					checkGet(t, c, key, key+10, true)
				}
			}
			c.Set(tc.key, 99)
			checkGet(t, c, tc.key, 99, true)
		})
	}
}

func TestClear(t *testing.T) {
	c := newTestCache[int, int](t, 2)
	c.Clear()
	c.Set(1, 10)
	c.Set(2, 20)
	c.Clear()
	c.Clear()
	if c.Len() != 0 {
		t.Fatalf("Len() = %d; want 0", c.Len())
	}
	checkGet(t, c, 1, 0, false)
	checkGet(t, c, 2, 0, false)

	// после очистки кеш снова заполняется до прежней вместимости
	c.Set(1, 100)
	c.Set(2, 200)
	c.Set(3, 300)
	checkGet(t, c, 1, 0, false)
	checkGet(t, c, 2, 200, true)
	checkGet(t, c, 3, 300, true)
	if c.Len() != 2 {
		t.Fatalf("Len() = %d; want 2", c.Len())
	}
}

func TestZeroValue(t *testing.T) {
	c := newTestCache[string, *int](t, 2)
	c.Set("", nil)
	// сохраненный nil отличается от отсутствующего ключа
	checkGet(t, c, "", (*int)(nil), true)
	checkGet(t, c, "missing", (*int)(nil), false)
}

func TestConcurrentAccess(t *testing.T) {
	c := newTestCache[int, int](t, 8)
	var wg sync.WaitGroup
	for worker := range 8 {
		wg.Go(func() {
			for i := range 100 {
				key := (worker + i) % 16
				c.Set(key, key)
				if got, ok := c.Get(key); ok && got != key {
					t.Errorf("Get(%d) = %d; want %d", key, got, key)
				}
				c.Remove(key)
				if size := c.Len(); size > 8 {
					t.Errorf("Len() = %d; want <= 8", size)
				}
				if i%10 == 0 {
					c.Clear()
				}
			}
		})
	}
	wg.Wait()
}

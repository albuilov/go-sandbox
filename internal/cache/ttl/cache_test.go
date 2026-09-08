package ttl

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// время двигаем вручную, без ожиданий и гонок между горутинами
func newTestCache[K comparable, V any](t *testing.T) (*ttlCache[K, V], func(time.Duration)) {
	t.Helper()
	cache, err := NewCache[K, V](5 * time.Second)
	if err != nil {
		t.Fatal(err)
	}
	c := cache.(*ttlCache[K, V])
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	var elapsed atomic.Int64
	c.now = func() time.Time { return start.Add(time.Duration(elapsed.Load())) }
	return c, func(d time.Duration) { elapsed.Add(int64(d)) }
}

func checkGet[K comparable, V comparable](t *testing.T, c Cache[K, V], key K, want V, exists bool) {
	t.Helper()
	if got, ok := c.Get(key); got != want || ok != exists {
		t.Fatalf("Get(%v) = (%v, %v); want (%v, %v)", key, got, ok, want, exists)
	}
}

func TestNewCache(t *testing.T) {
	for _, ttl := range []time.Duration{-time.Second, 0, time.Nanosecond, time.Second} {
		t.Run(ttl.String(), func(t *testing.T) {
			c, err := NewCache[int, string](ttl)
			if ttl <= 0 {
				if err == nil || c != nil {
					t.Fatalf("NewCache(%s) = (%v, %v); want nil and error", ttl, c, err)
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

func TestExpiration(t *testing.T) {
	for _, elapsed := range []time.Duration{0, 5*time.Second - time.Nanosecond, 5 * time.Second, 6 * time.Second} {
		t.Run(elapsed.String(), func(t *testing.T) {
			c, advance := newTestCache[string, int](t)
			c.Set("A", 10)
			advance(elapsed)
			want, exists := 10, elapsed < 5*time.Second
			if !exists {
				want = 0
			}
			checkGet(t, c, "A", want, exists)
			// проверяем удаление напрямую, Len сам мог бы очистить запись
			if _, ok := c.nodes["A"]; ok != exists {
				t.Fatalf("stored A = %v; want %v", ok, exists)
			}
		})
	}
}

func TestGetDoesNotExtendTTL(t *testing.T) {
	c, advance := newTestCache[string, int](t)
	c.Set("A", 10)
	advance(4 * time.Second)
	checkGet(t, c, "A", 10, true)
	advance(time.Second)
	checkGet(t, c, "A", 0, false)
}

func TestSetRenewsTTL(t *testing.T) {
	for _, elapsed := range []time.Duration{time.Second, 5 * time.Second, 6 * time.Second} {
		t.Run(elapsed.String(), func(t *testing.T) {
			c, advance := newTestCache[string, int](t)
			c.Set("A", 10)
			advance(elapsed)
			// заменяем и действующую, и уже просроченную запись
			c.Set("A", 20)
			if c.Len() != 1 {
				t.Fatalf("Len() = %d; want 1", c.Len())
			}
			advance(5*time.Second - time.Nanosecond)
			checkGet(t, c, "A", 20, true)
			advance(time.Nanosecond)
			checkGet(t, c, "A", 0, false)
		})
	}
}

func TestCleanup(t *testing.T) {
	for _, operation := range []struct {
		name string
		run  func(Cache[string, int]) int
	}{
		{"DeleteExpired", func(c Cache[string, int]) int { c.DeleteExpired(); return -1 }},
		{"Len", func(c Cache[string, int]) int { return c.Len() }},
	} {
		t.Run(operation.name, func(t *testing.T) {
			c, advance := newTestCache[string, int](t)
			operation.run(c)
			c.Set("A", 10)
			advance(2 * time.Second)
			c.Set("B", 20)
			for _, step := range []struct {
				elapsed time.Duration
				want    int
			}{{0, 2}, {3 * time.Second, 1}, {2 * time.Second, 0}, {0, 0}} {
				advance(step.elapsed)
				got := operation.run(c)
				if operation.name == "Len" && got != step.want {
					t.Fatalf("Len() = %d; want %d", got, step.want)
				}
				if len(c.nodes) != step.want {
					t.Fatalf("stored count = %d; want %d", len(c.nodes), step.want)
				}
				if step.want == 1 {
					checkGet(t, c, "A", 0, false)
					checkGet(t, c, "B", 20, true)
				}
			}
		})
	}
}

func TestRemove(t *testing.T) {
	for _, elapsed := range []time.Duration{0, 5 * time.Second} {
		t.Run(elapsed.String(), func(t *testing.T) {
			c, advance := newTestCache[int, int](t)
			c.Remove(0)
			c.Set(1, 10)
			advance(elapsed)
			c.Set(2, 20)
			c.Remove(1)
			c.Remove(1)
			if _, ok := c.nodes[1]; ok {
				t.Fatal("removed key still stored")
			}
			checkGet(t, c, 1, 0, false)
			checkGet(t, c, 2, 20, true)
			c.Set(1, 30)
			checkGet(t, c, 1, 30, true)
		})
	}
}

func TestClear(t *testing.T) {
	c, advance := newTestCache[int, int](t)
	c.Clear()
	c.Set(1, 10)
	advance(5 * time.Second)
	c.Set(2, 20)
	c.Clear()
	c.Clear()
	if c.Len() != 0 {
		t.Fatalf("Len() = %d; want 0", c.Len())
	}
	checkGet(t, c, 1, 0, false)
	checkGet(t, c, 2, 0, false)
	// после очистки сохраняется прежний TTL
	c.Set(1, 30)
	advance(4 * time.Second)
	checkGet(t, c, 1, 30, true)
	advance(time.Second)
	checkGet(t, c, 1, 0, false)
}

func TestZeroValue(t *testing.T) {
	c, advance := newTestCache[string, *int](t)
	c.Set("", nil)
	checkGet(t, c, "", (*int)(nil), true)
	checkGet(t, c, "missing", (*int)(nil), false)
	advance(5 * time.Second)
	checkGet(t, c, "", (*int)(nil), false)
}

func TestConcurrentAccess(t *testing.T) {
	c, advance := newTestCache[int, int](t)
	var wg sync.WaitGroup
	for worker := range 8 {
		wg.Go(func() {
			for i := range 100 {
				key := (worker + i) % 16
				c.Set(key, key)
				if got, ok := c.Get(key); ok && got != key {
					t.Errorf("Get(%d) = %d; want %d", key, got, key)
				}
				advance(time.Second)
				c.DeleteExpired()
				c.Remove(key)
				if size := c.Len(); size > 16 {
					t.Errorf("Len() = %d; want <= 16", size)
				}
				if i%10 == 0 {
					c.Clear()
				}
			}
		})
	}
	wg.Wait()
	advance(5 * time.Second)
	if c.Len() != 0 {
		t.Fatalf("Len() = %d; want 0", c.Len())
	}
}

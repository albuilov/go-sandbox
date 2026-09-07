package weighted

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"
)

func newTestSemaphore(t *testing.T, limit int) *Semaphore {
	t.Helper()

	s, err := NewSemaphore(limit)
	if err != nil {
		t.Fatal(err)
	}

	return s
}

func TestNewSemaphore(t *testing.T) {
	for _, limit := range []int{-1, 0} {
		if s, err := NewSemaphore(limit); err == nil || s != nil {
			t.Errorf("NewSemaphore(%d) = %v, %v; want nil, error", limit, s, err)
		}
	}
}

func TestInvalidWeight(t *testing.T) {
	for _, weight := range []int{-1, 0, 6} {
		s := newTestSemaphore(t, 5)

		if err := s.Acquire(context.Background(), weight); err == nil {
			t.Errorf("Acquire(%d) accepted invalid weight", weight)
		}

		if s.TryAcquire(weight) {
			t.Errorf("TryAcquire(%d) accepted invalid weight", weight)
		}

		if err := s.Release(weight); err == nil {
			t.Errorf("Release(%d) accepted invalid weight", weight)
		}

		if !s.TryAcquire(5) {
			t.Fatal("invalid operation changed capacity")
		}

		if err := s.Release(5); err != nil {
			t.Fatal(err)
		}
	}
}

func TestTryAcquireCapacity(t *testing.T) {
	s := newTestSemaphore(t, 5)

	if !s.TryAcquire(3) {
		t.Fatal("TryAcquire(3) failed")
	}

	if s.TryAcquire(3) {
		t.Fatal("TryAcquire exceeded limit: 3 + 3 > 5")
	}

	if !s.TryAcquire(2) {
		t.Fatal("failed TryAcquire changed available weight")
	}

	if s.TryAcquire(1) {
		t.Fatal("TryAcquire succeeded on full semaphore")
	}

	if err := s.Release(5); err != nil {
		t.Fatal(err)
	}

	if !s.TryAcquire(5) {
		t.Fatal("released weight unavailable")
	}

	if err := s.Release(5); err != nil {
		t.Fatal(err)
	}
}

func TestReleaseTooMuch(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newTestSemaphore(t, 5)

		if !s.TryAcquire(2) {
			t.Fatal("initial acquire failed")
		}

		if err := s.Release(3); err == nil {
			t.Fatal("Release accepted more weight than held")
		}

		// ошибка не должна освободить уже занятые места
		done := make(chan error, 1)
		go func() {
			done <- s.Acquire(context.Background(), 4)
		}()

		synctest.Wait()

		select {
		case err := <-done:
			t.Fatalf("invalid Release changed capacity: %v", err)
		default:
		}

		if err := s.Release(2); err != nil {
			t.Fatal(err)
		}

		if err := <-done; err != nil {
			t.Fatal(err)
		}

		if err := s.Release(4); err != nil {
			t.Fatal(err)
		}
	})
}

func TestAcquireCancelledContext(t *testing.T) {
	s := newTestSemaphore(t, 5)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan error, 1)

	go func() {
		done <- s.Acquire(ctx, 3)
	}()

	// реальный таймаут защищает тест от повторного Lock внутри Acquire
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Acquire = %v; want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Acquire hung with already cancelled context")
	}

	if !s.TryAcquire(5) {
		t.Fatal("cancelled Acquire took permits")
	}

	if err := s.Release(5); err != nil {
		t.Fatal(err)
	}
}

func TestAcquireWaitsForWholeWeight(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newTestSemaphore(t, 5)

		if err := s.Acquire(context.Background(), 4); err != nil {
			t.Fatal(err)
		}

		done := make(chan error, 1)

		go func() {
			done <- s.Acquire(context.Background(), 3)
		}()

		synctest.Wait()

		select {
		case err := <-done:
			t.Fatalf("Acquire returned without enough capacity: %v", err)
		default:
		}

		// одного освобожденного места еще недостаточно
		if err := s.Release(1); err != nil {
			t.Fatal(err)
		}

		synctest.Wait()
		select {
		case err := <-done:
			t.Fatalf("Acquire accepted partial capacity: %v", err)
		default:
		}

		if err := s.Release(1); err != nil {
			t.Fatal(err)
		}

		synctest.Wait()

		select {
		case err := <-done:
			if err != nil {
				t.Fatal(err)
			}
		default:
			t.Fatal("Acquire did not resume")
		}

		if err := s.Release(5); err != nil {
			t.Fatal(err)
		}
	})
}

func TestCancelWaitingAcquire(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := newTestSemaphore(t, 5)

		if err := s.Acquire(context.Background(), 5); err != nil {
			t.Fatal(err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		done := make(chan error, 1)

		go func() {
			done <- s.Acquire(ctx, 3)
		}()

		synctest.Wait()
		cancel()

		if err := <-done; !errors.Is(err, context.Canceled) {
			t.Fatalf("Acquire = %v; want context.Canceled", err)
		}

		if err := s.Release(5); err != nil {
			t.Fatal(err)
		}

		if err := s.Acquire(context.Background(), 5); err != nil {
			t.Fatal(err)
		}

		if err := s.Release(5); err != nil {
			t.Fatal(err)
		}
	})
}

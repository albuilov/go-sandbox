package cancellable

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"
)

func TestNewSemaphore(t *testing.T) {
	for _, limit := range []int{-1, 0} {
		if s, err := NewSemaphore(limit); err == nil || s != nil {
			t.Errorf("NewSemaphore(%d) = %v, %v; want nil, error", limit, s, err)
		}
	}
}

func TestTryAcquire(t *testing.T) {
	s, err := NewSemaphore(2)
	if err != nil {
		t.Fatal(err)
	}

	if !s.TryAcquire() || !s.TryAcquire() {
		t.Fatal("could not fill capacity")
	}

	if s.TryAcquire() {
		t.Fatal("TryAcquire exceeded limit")
	}

	if err := s.Release(); err != nil {
		t.Fatal(err)
	}

	if !s.TryAcquire() {
		t.Fatal("released permit unavailable")
	}

	if err := s.Release(); err != nil {
		t.Fatal(err)
	}
	if err := s.Release(); err != nil {
		t.Fatal(err)
	}
}

func TestAcquireCancelledContext(t *testing.T) {
	s, err := NewSemaphore(1)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := s.Acquire(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Acquire = %v; want context.Canceled", err)
	}

	if !s.TryAcquire() {
		t.Fatal("cancelled Acquire took a permit")
	}

	_ = s.Release()
}

func TestAcquireWaiting(t *testing.T) {
	for _, scenario := range []string{"release", "cancel", "deadline"} {
		t.Run(scenario, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s, err := NewSemaphore(1)
				if err != nil {
					t.Fatal(err)
				}

				if !s.TryAcquire() {
					t.Fatal("initial acquire failed")
				}

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				done := make(chan error, 1)
				go func() {
					done <- s.Acquire(ctx)
				}()

				synctest.Wait()

				select {
				case err := <-done:
					t.Fatalf("Acquire returned while full: %v", err)
				default:
				}

				switch scenario {
				case "release":
					s.Release()
				case "cancel":
					cancel()
				case "deadline":
					time.Sleep(time.Second) // В synctest время виртуальное.
				}

				err = <-done

				var want error
				if scenario == "cancel" {
					want = context.Canceled
				}

				if scenario == "deadline" {
					want = context.DeadlineExceeded
				}

				if !errors.Is(err, want) {
					t.Fatalf("Acquire = %v; want %v", err, want)
				}

				// При успехе освобождаем новое разрешение, при отмене — исходное.
				if err := s.Release(); err != nil {
					t.Fatal(err)
				}

				if !s.TryAcquire() {
					t.Fatal("permit lost")
				}

				if err := s.Release(); err != nil {
					t.Fatal(err)
				}
			})
		})
	}
}

func TestReleaseWithoutAcquire(t *testing.T) {
	s, err := NewSemaphore(1)
	if err != nil {
		t.Fatal(err)
	}

	if err := s.Release(); err == nil {
		t.Error("Release without Acquire must panic")
	}
}

package basic

import (
	"testing"
	"testing/synctest"
)

func TestNewSemaphore(t *testing.T) {
	for _, limit := range []int{-1, 0} {
		if s, err := NewSemaphore(limit); err == nil || s != nil {
			t.Errorf("NewSemaphore(%d) = %v, %v; want nil, error", limit, s, err)
		}
	}
}

func TestAcquireWaitsForRelease(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, err := NewSemaphore(2)
		if err != nil {
			t.Fatal(err)
		}

		s.Acquire()
		s.Acquire()

		acquired := make(chan struct{})

		go func() {
			s.Acquire()
			close(acquired)
		}()

		// Даем горутине дойти до ожидания места
		synctest.Wait()

		select {
		case <-acquired:
			t.Fatal("Acquire exceeded limit")
		default:
		}
		s.Release()

		synctest.Wait()

		select {
		case <-acquired:
		default:
			t.Fatal("Acquire did not resume after Release")
		}

		s.Release()
		s.Release()

		// Все места можно занять снова
		s.Acquire()
		s.Acquire()

		s.Release()
		s.Release()
	})
}

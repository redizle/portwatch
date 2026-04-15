package throttle_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

// TestConcurrentAllow verifies that Allow is safe under concurrent access
// and that exactly one goroutine wins per interval window.
func TestConcurrentAllow(t *testing.T) {
	th := throttle.New(50 * time.Millisecond)

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		allowed int
	)

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if th.Allow() {
				mu.Lock()
				allowed++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if allowed != 1 {
		t.Fatalf("expected exactly 1 allowed call, got %d", allowed)
	}

	skipped := th.Skipped()
	if skipped != int64(20-allowed) {
		t.Fatalf("expected %d skipped, got %d", 20-allowed, skipped)
	}
}

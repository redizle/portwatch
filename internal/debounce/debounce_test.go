package debounce_test

import (
	"sync"
	"testing"
	"time"

	"portwatch/internal/debounce"
)

func TestPush_FiresAfterWindow(t *testing.T) {
	var mu sync.Mutex
	fired := []string{}

	d := debounce.New(30*time.Millisecond, func(key string) {
		mu.Lock()
		fired = append(fired, key)
		mu.Unlock()
	})

	d.Push("8080")
	time.Sleep(60 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(fired) != 1 || fired[0] != "8080" {
		t.Fatalf("expected fired=[8080], got %v", fired)
	}
}

func TestPush_ResetsTimer(t *testing.T) {
	var mu sync.Mutex
	count := 0

	d := debounce.New(40*time.Millisecond, func(_ string) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	d.Push("443")
	time.Sleep(20 * time.Millisecond)
	d.Push("443") // reset
	time.Sleep(20 * time.Millisecond)
	d.Push("443") // reset again
	time.Sleep(70 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 1 {
		t.Fatalf("expected handler called once, got %d", count)
	}
}

func TestCancel_Preventsfire(t *testing.T) {
	fired := false
	d := debounce.New(30*time.Millisecond, func(_ string) { fired = true })

	d.Push("22")
	d.Cancel("22")
	time.Sleep(60 * time.Millisecond)

	if fired {
		t.Fatal("expected handler not to fire after Cancel")
	}
}

func TestPending_Count(t *testing.T) {
	d := debounce.New(200*time.Millisecond, func(_ string) {})

	d.Push("80")
	d.Push("443")

	if p := d.Pending(); p != 2 {
		t.Fatalf("expected 2 pending, got %d", p)
	}

	d.Cancel("80")
	if p := d.Pending(); p != 1 {
		t.Fatalf("expected 1 pending after cancel, got %d", p)
	}
}

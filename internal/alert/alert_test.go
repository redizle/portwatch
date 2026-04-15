package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func TestSend_NoHooks(t *testing.T) {
	n := alert.New(nil, 0)
	evt := alert.Event{Port: 8080, State: "open", Timestamp: time.Now(), Host: "localhost"}
	if err := n.Send(evt); err != nil {
		t.Fatalf("expected no error with no hooks, got: %v", err)
	}
}

func TestSend_SuccessfulHook(t *testing.T) {
	var received atomic.Int32
	var capturedEvent alert.Event

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&capturedEvent); err != nil {
			t.Errorf("decode body: %v", err)
		}
		received.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := alert.New([]string{ts.URL}, 3*time.Second)
	evt := alert.Event{Port: 9090, State: "closed", Timestamp: time.Now(), Host: "localhost"}

	if err := n.Send(evt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Load() != 1 {
		t.Errorf("expected 1 request, got %d", received.Load())
	}
	if capturedEvent.Port != 9090 || capturedEvent.State != "closed" {
		t.Errorf("unexpected event payload: %+v", capturedEvent)
	}
}

func TestSend_HookReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := alert.New([]string{ts.URL}, 3*time.Second)
	evt := alert.Event{Port: 443, State: "open", Timestamp: time.Now(), Host: "localhost"}

	if err := n.Send(evt); err == nil {
		t.Fatal("expected error for non-2xx response, got nil")
	}
}

func TestSend_UnreachableHook(t *testing.T) {
	n := alert.New([]string{"http://127.0.0.1:1"}, 500*time.Millisecond)
	evt := alert.Event{Port: 22, State: "open", Timestamp: time.Now(), Host: "localhost"}

	if err := n.Send(evt); err == nil {
		t.Fatal("expected error for unreachable hook, got nil")
	}
}

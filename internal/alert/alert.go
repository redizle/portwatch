package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Event represents a port state change that triggered an alert.
type Event struct {
	Port      int       `json:"port"`
	State     string    `json:"state"` // "open" or "closed"
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
}

// Hook defines an alerting destination.
type Hook struct {
	URL     string
	Timeout time.Duration
}

// Notifier sends alerts to configured hooks.
type Notifier struct {
	hooks  []Hook
	client *http.Client
}

// New creates a Notifier with the given webhook URLs.
func New(urls []string, timeout time.Duration) *Notifier {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	hooks := make([]Hook, 0, len(urls))
	for _, u := range urls {
		hooks = append(hooks, Hook{URL: u, Timeout: timeout})
	}
	return &Notifier{
		hooks:  hooks,
		client: &http.Client{Timeout: timeout},
	}
}

// Send dispatches the event to all configured hooks.
// It returns a combined error if any hook fails.
func (n *Notifier) Send(evt Event) error {
	if len(n.hooks) == 0 {
		return nil
	}

	body, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("alert: marshal event: %w", err)
	}

	var errs []error
	for _, h := range n.hooks {
		resp, err := n.client.Post(h.URL, "application/json", bytes.NewReader(body))
		if err != nil {
			errs = append(errs, fmt.Errorf("hook %s: %w", h.URL, err))
			continue
		}
		resp.Body.Close()
		if resp.StatusCode >= 300 {
			errs = append(errs, fmt.Errorf("hook %s: unexpected status %d", h.URL, resp.StatusCode))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("alert: %d hook(s) failed: %v", len(errs), errs)
	}
	return nil
}

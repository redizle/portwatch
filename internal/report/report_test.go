package report_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/report"
	"github.com/user/portwatch/internal/state"
)

func newTestReporter(t *testing.T) (*report.Reporter, *history.History, *state.State, *bytes.Buffer) {
	t.Helper()
	h := history.New(50, "")
	s := state.New()
	var buf bytes.Buffer
	r := report.New(h, s, &buf)
	return r, h, s, &buf
}

func TestSummary_NoPorts(t *testing.T) {
	r, _, _, buf := newTestReporter(t)
	r.Summary()
	if !strings.Contains(buf.String(), "No ports") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestSummary_WithPorts(t *testing.T) {
	r, _, s, buf := newTestReporter(t)
	s.Update(8080, true)
	s.Update(9090, false)
	r.Summary()
	out := buf.String()
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
	if !strings.Contains(out, "open") {
		t.Errorf("expected 'open' status in output")
	}
	if !strings.Contains(out, "closed") {
		t.Errorf("expected 'closed' status in output")
	}
}

func TestRecentActivity_NoEntries(t *testing.T) {
	r, _, _, buf := newTestReporter(t)
	r.RecentActivity(10)
	if !strings.Contains(buf.String(), "No recent activity") {
		t.Errorf("expected empty activity message, got: %s", buf.String())
	}
}

func TestRecentActivity_WithEntries(t *testing.T) {
	r, h, _, buf := newTestReporter(t)
	h.Record(8080, "opened", time.Now())
	h.Record(443, "closed", time.Now())
	r.RecentActivity(5)
	out := buf.String()
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in activity output")
	}
	if !strings.Contains(out, "opened") {
		t.Errorf("expected event 'opened' in activity output")
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	h := history.New(50, "")
	s := state.New()
	// should not panic
	rep := report.New(h, s, nil)
	if rep == nil {
		t.Fatal("expected non-nil reporter")
	}
}

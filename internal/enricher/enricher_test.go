package enricher_test

import (
	"testing"

	"github.com/user/portwatch/internal/enricher"
	"github.com/user/portwatch/internal/portname"
	"github.com/user/portwatch/internal/tagger"
)

func TestEnrich_WellKnownPort(t *testing.T) {
	e := enricher.New(nil, nil)
	m := e.Enrich(80)
	if m.Port != 80 {
		t.Fatalf("expected port 80, got %d", m.Port)
	}
	if m.Name != "http" {
		t.Errorf("expected name http, got %q", m.Name)
	}
}

func TestEnrich_UnknownPort(t *testing.T) {
	e := enricher.New(nil, nil)
	m := e.Enrich(19999)
	if m.Name == "" {
		t.Error("expected non-empty fallback name")
	}
}

func TestEnrich_WithOverrides(t *testing.T) {
	r := portname.New(map[int]string{9100: "metrics"})
	tg := tagger.New(map[int]string{9100: "observability"})
	e := enricher.New(r, tg)

	m := e.Enrich(9100)
	if m.Name != "metrics" {
		t.Errorf("expected metrics, got %q", m.Name)
	}
	if m.Tag != "observability" {
		t.Errorf("expected observability tag, got %q", m.Tag)
	}
}

func TestEnrichAll_ReturnsAllPorts(t *testing.T) {
	e := enricher.New(nil, nil)
	ports := []int{22, 80, 443}
	results := e.EnrichAll(ports)
	if len(results) != len(ports) {
		t.Fatalf("expected %d results, got %d", len(ports), len(results))
	}
	for i, m := range results {
		if m.Port != ports[i] {
			t.Errorf("index %d: expected port %d, got %d", i, ports[i], m.Port)
		}
	}
}

func TestEnrichAll_EmptySlice(t *testing.T) {
	e := enricher.New(nil, nil)
	results := e.EnrichAll(nil)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

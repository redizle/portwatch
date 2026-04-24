package portprofile

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	p := Profile{Name: "web", Label: "HTTP", Owner: "alice", Priority: 1, Tags: []string{"public"}}
	if err := s.Set(80, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := s.Get(80)
	if !ok {
		t.Fatal("expected profile to be found")
	}
	if got.Name != "web" {
		t.Errorf("expected name 'web', got %q", got.Name)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	s := New()
	if err := s.Set(0, Profile{Name: "x"}); err == nil {
		t.Error("expected error for port 0")
	}
	if err := s.Set(65536, Profile{Name: "x"}); err == nil {
		t.Error("expected error for port 65536")
	}
}

func TestSet_EmptyName(t *testing.T) {
	s := New()
	if err := s.Set(443, Profile{Name: ""}); err == nil {
		t.Error("expected error for empty profile name")
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	_, ok := s.Get(9999)
	if ok {
		t.Error("expected missing profile")
	}
}

func TestRemove(t *testing.T) {
	s := New()
	_ = s.Set(22, Profile{Name: "ssh"})
	s.Remove(22)
	_, ok := s.Get(22)
	if ok {
		t.Error("expected profile to be removed")
	}
}

func TestRemove_NotPresent(t *testing.T) {
	s := New()
	// should not panic
	s.Remove(1234)
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := New()
	_ = s.Set(80, Profile{Name: "web"})
	_ = s.Set(443, Profile{Name: "tls"})
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
	// mutating the copy must not affect the store
	delete(all, 80)
	if _, ok := s.Get(80); !ok {
		t.Error("store was mutated via All() copy")
	}
}

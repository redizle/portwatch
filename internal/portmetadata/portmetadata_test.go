package portmetadata

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	if err := s.Set(8080, "env", "production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := s.Get(8080, "env")
	if !ok || v != "production" {
		t.Fatalf("expected production, got %q ok=%v", v, ok)
	}
}

func TestSet_InvalidPort(t *testing.T) {
	s := New()
	if err := s.Set(0, "k", "v"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := s.Set(70000, "k", "v"); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestSet_EmptyKey(t *testing.T) {
	s := New()
	if err := s.Set(443, "", "value"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	_, ok := s.Get(9090, "missing")
	if ok {
		t.Fatal("expected false for missing entry")
	}
}

func TestDelete_RemovesKey(t *testing.T) {
	s := New()
	_ = s.Set(22, "owner", "alice")
	s.Delete(22, "owner")
	_, ok := s.Get(22, "owner")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestDelete_CleansUpEmptyPort(t *testing.T) {
	s := New()
	_ = s.Set(22, "only", "value")
	s.Delete(22, "only")
	if len(s.All(22)) != 0 {
		t.Fatal("expected port entry to be removed after last key deleted")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := New()
	_ = s.Set(80, "a", "1")
	_ = s.Set(80, "b", "2")
	all := s.All(80)
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	all["c"] = "3"
	if _, ok := s.Get(80, "c"); ok {
		t.Fatal("mutation of copy should not affect store")
	}
}

func TestClear_RemovesPort(t *testing.T) {
	s := New()
	_ = s.Set(3306, "db", "mysql")
	s.Clear(3306)
	if len(s.All(3306)) != 0 {
		t.Fatal("expected all metadata cleared")
	}
}

func TestSet_MultipleKeys(t *testing.T) {
	s := New()
	_ = s.Set(5432, "db", "postgres")
	_ = s.Set(5432, "env", "staging")
	all := s.All(5432)
	if all["db"] != "postgres" || all["env"] != "staging" {
		t.Fatalf("unexpected values: %v", all)
	}
}

package porttag

import (
	"testing"
)

func TestAdd_And_Get(t *testing.T) {
	s := New()
	if err := s.Add(8080, "web"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags := s.Get(8080)
	if len(tags) != 1 || tags[0] != "web" {
		t.Fatalf("expected [web], got %v", tags)
	}
}

func TestAdd_InvalidPort(t *testing.T) {
	s := New()
	if err := s.Add(0, "web"); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestAdd_EmptyTag(t *testing.T) {
	s := New()
	if err := s.Add(80, ""); err == nil {
		t.Fatal("expected error for empty tag")
	}
}

func TestHas_True(t *testing.T) {
	s := New()
	_ = s.Add(443, "tls")
	if !s.Has(443, "tls") {
		t.Fatal("expected Has to return true")
	}
}

func TestHas_False(t *testing.T) {
	s := New()
	if s.Has(443, "tls") {
		t.Fatal("expected Has to return false")
	}
}

func TestRemove(t *testing.T) {
	s := New()
	_ = s.Add(22, "ssh")
	_ = s.Add(22, "secure")
	s.Remove(22, "ssh")
	if s.Has(22, "ssh") {
		t.Fatal("expected ssh tag to be removed")
	}
	if !s.Has(22, "secure") {
		t.Fatal("expected secure tag to remain")
	}
}

func TestRemove_LastTag_ClearsPort(t *testing.T) {
	s := New()
	_ = s.Add(9000, "custom")
	s.Remove(9000, "custom")
	if len(s.Get(9000)) != 0 {
		t.Fatal("expected no tags after removing last")
	}
}

func TestClear(t *testing.T) {
	s := New()
	_ = s.Add(3306, "db")
	_ = s.Add(3306, "mysql")
	s.Clear(3306)
	if len(s.Get(3306)) != 0 {
		t.Fatal("expected all tags cleared")
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	if tags := s.Get(1234); len(tags) != 0 {
		t.Fatalf("expected empty slice, got %v", tags)
	}
}

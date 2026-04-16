package baseline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestAdd_Contains(t *testing.T) {
	b := New()
	b.Add(8080, "dev server")
	if !b.Contains(8080) {
		t.Fatal("expected 8080 to be in baseline")
	}
}

func TestRemove(t *testing.T) {
	b := New()
	b.Add(8080, "")
	b.Remove(8080)
	if b.Contains(8080) {
		t.Fatal("expected 8080 to be removed")
	}
}

func TestContains_Missing(t *testing.T) {
	b := New()
	if b.Contains(9999) {
		t.Fatal("expected false for missing port")
	}
}

func TestUnexpected(t *testing.T) {
	b := New()
	b.Add(80, "http")
	b.Add(443, "https")
	active := []int{80, 443, 8080, 3306}
	unexpected := b.Unexpected(active)
	if len(unexpected) != 2 {
		t.Fatalf("expected 2 unexpected ports, got %d", len(unexpected))
	}
}

func TestUnexpected_AllBaselined(t *testing.T) {
	b := New()
	b.Add(80, "")
	result := b.Unexpected([]int{80})
	if len(result) != 0 {
		t.Fatal("expected no unexpected ports")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	b := New()
	b.Add(22, "ssh")
	b.Add(80, "http")
	all := b.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestSave_AndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	b := New()
	b.Add(80, "http")
	b.Add(22, "ssh")

	if err := b.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !loaded.Contains(80) || !loaded.Contains(22) {
		t.Fatal("loaded baseline missing expected ports")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0644)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSave_CreatesValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	b := New()
	b.Add(443, "https")
	_ = b.Save(path)
	data, _ := os.ReadFile(path)
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("saved file is not valid JSON: %v", err)
	}
}

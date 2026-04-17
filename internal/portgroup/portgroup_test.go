package portgroup_test

import (
	"testing"

	"portwatch/internal/portgroup"
)

func TestAdd_ValidGroup(t *testing.T) {
	r := portgroup.New()
	if err := r.Add("web", []int{80, 443}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g, ok := r.Get("web")
	if !ok {
		t.Fatal("expected group to exist")
	}
	if len(g.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(g.Ports))
	}
}

func TestAdd_InvalidPort(t *testing.T) {
	r := portgroup.New()
	if err := r.Add("bad", []int{0}); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := r.Add("bad", []int{65536}); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestGet_Missing(t *testing.T) {
	r := portgroup.New()
	_, ok := r.Get("nope")
	if ok {
		t.Fatal("expected false for missing group")
	}
}

func TestRemove(t *testing.T) {
	r := portgroup.New()
	_ = r.Add("db", []int{5432})
	r.Remove("db")
	_, ok := r.Get("db")
	if ok {
		t.Fatal("expected group to be removed")
	}
}

func TestGroupsFor_ReturnsMatchingGroups(t *testing.T) {
	r := portgroup.New()
	_ = r.Add("web", []int{80, 443})
	_ = r.Add("all", []int{80, 22, 443})
	names := r.GroupsFor(80)
	if len(names) != 2 {
		t.Fatalf("expected 2 groups for port 80, got %d", len(names))
	}
}

func TestGroupsFor_NoMatch(t *testing.T) {
	r := portgroup.New()
	_ = r.Add("web", []int{80})
	names := r.GroupsFor(9999)
	if len(names) != 0 {
		t.Fatalf("expected no groups, got %d", len(names))
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	r := portgroup.New()
	_ = r.Add("a", []int{1111})
	_ = r.Add("b", []int{2222})
	all := r.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(all))
	}
}

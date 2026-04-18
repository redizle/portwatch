package portlabel_test

import (
	"testing"

	"portwatch/internal/portlabel"
)

func TestSet_And_Get(t *testing.T) {
	l := portlabel.New(nil)
	lb := portlabel.Label{Port: 8080, Name: "http-alt", Description: "Alt HTTP", Color: "green"}
	if err := l.Set(lb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := l.Get(8080)
	if !ok {
		t.Fatal("expected label to exist")
	}
	if got.Name != "http-alt" {
		t.Errorf("got name %q, want %q", got.Name, "http-alt")
	}
}

func TestSet_InvalidPort(t *testing.T) {
	l := portlabel.New(nil)
	err := l.Set(portlabel.Label{Port: 0, Name: "bad"})
	if err != portlabel.ErrInvalidPort {
		t.Errorf("expected ErrInvalidPort, got %v", err)
	}
}

func TestGet_Missing(t *testing.T) {
	l := portlabel.New(nil)
	_, ok := l.Get(9999)
	if ok {
		t.Error("expected missing label")
	}
}

func TestRemove(t *testing.T) {
	l := portlabel.New(nil)
	_ = l.Set(portlabel.Label{Port: 443, Name: "https"})
	l.Remove(443)
	_, ok := l.Get(443)
	if ok {
		t.Error("expected label to be removed")
	}
}

func TestNew_Seed(t *testing.T) {
	seed := []portlabel.Label{
		{Port: 22, Name: "ssh", Color: "yellow"},
		{Port: 80, Name: "http", Color: "green"},
	}
	l := portlabel.New(seed)
	if _, ok := l.Get(22); !ok {
		t.Error("expected seed label for port 22")
	}
	if _, ok := l.Get(80); !ok {
		t.Error("expected seed label for port 80")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	l := portlabel.New(nil)
	_ = l.Set(portlabel.Label{Port: 3306, Name: "mysql"})
	_ = l.Set(portlabel.Label{Port: 5432, Name: "postgres"})
	all := l.All()
	if len(all) != 2 {
		t.Errorf("expected 2 labels, got %d", len(all))
	}
}

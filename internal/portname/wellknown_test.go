package portname

import "testing"

func TestWellKnown_ContainsHTTP(t *testing.T) {
	name, ok := wellKnown[80]
	if !ok {
		t.Fatal("expected port 80 in wellKnown map")
	}
	if name != "http" {
		t.Errorf("expected http, got %s", name)
	}
}

func TestWellKnown_ContainsSSH(t *testing.T) {
	name, ok := wellKnown[22]
	if !ok {
		t.Fatal("expected port 22 in wellKnown map")
	}
	if name != "ssh" {
		t.Errorf("expected ssh, got %s", name)
	}
}

func TestWellKnown_ContainsRedis(t *testing.T) {
	name, ok := wellKnown[6379]
	if !ok {
		t.Fatal("expected port 6379 in wellKnown map")
	}
	if name != "redis" {
		t.Errorf("expected redis, got %s", name)
	}
}

func TestWellKnown_UnknownPort(t *testing.T) {
	_, ok := wellKnown[9999]
	if ok {
		t.Error("expected port 9999 to be absent from wellKnown map")
	}
}

func TestWellKnown_AllEntriesNonEmpty(t *testing.T) {
	for port, name := range wellKnown {
		if name == "" {
			t.Errorf("port %d has empty name", port)
		}
		if port <= 0 || port > 65535 {
			t.Errorf("port %d is out of valid range", port)
		}
	}
}

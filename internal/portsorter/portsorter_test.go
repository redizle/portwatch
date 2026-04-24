package portsorter_test

import (
	"testing"

	"github.com/user/portwatch/internal/portsorter"
)

func sampleEntries() []portsorter.Entry {
	return []portsorter.Entry{
		{Port: 8080, Status: "open", Label: "http-alt"},
		{Port: 22, Status: "closed", Label: "ssh"},
		{Port: 443, Status: "open", Label: "https"},
		{Port: 3306, Status: "closed", Label: "mysql"},
	}
}

func TestSort_ByPort_Ascending(t *testing.T) {
	s := portsorter.New(portsorter.ByPort, portsorter.Ascending)
	entries := sampleEntries()
	s.Sort(entries)
	if entries[0].Port != 22 || entries[3].Port != 8080 {
		t.Errorf("expected ascending port order, got %v", entries)
	}
}

func TestSort_ByPort_Descending(t *testing.T) {
	s := portsorter.New(portsorter.ByPort, portsorter.Descending)
	entries := sampleEntries()
	s.Sort(entries)
	if entries[0].Port != 8080 || entries[3].Port != 22 {
		t.Errorf("expected descending port order, got %v", entries)
	}
}

func TestSort_ByStatus_Ascending(t *testing.T) {
	s := portsorter.New(portsorter.ByStatus, portsorter.Ascending)
	entries := sampleEntries()
	s.Sort(entries)
	// "closed" < "open" lexicographically
	if entries[0].Status != "closed" {
		t.Errorf("expected first entry to have status 'closed', got %q", entries[0].Status)
	}
}

func TestSort_ByLabel_Ascending(t *testing.T) {
	s := portsorter.New(portsorter.ByLabel, portsorter.Ascending)
	entries := sampleEntries()
	s.Sort(entries)
	// http-alt, https, mysql, ssh
	if entries[0].Label != "http-alt" {
		t.Errorf("expected first label 'http-alt', got %q", entries[0].Label)
	}
	if entries[3].Label != "ssh" {
		t.Errorf("expected last label 'ssh', got %q", entries[3].Label)
	}
}

func TestSortedCopy_DoesNotMutateOriginal(t *testing.T) {
	s := portsorter.New(portsorter.ByPort, portsorter.Ascending)
	original := sampleEntries()
	firstPort := original[0].Port
	_ = s.SortedCopy(original)
	if original[0].Port != firstPort {
		t.Error("SortedCopy mutated the original slice")
	}
}

func TestSortedCopy_ReturnsSortedResult(t *testing.T) {
	s := portsorter.New(portsorter.ByPort, portsorter.Ascending)
	result := s.SortedCopy(sampleEntries())
	for i := 1; i < len(result); i++ {
		if result[i].Port < result[i-1].Port {
			t.Errorf("result not sorted at index %d: %d < %d", i, result[i].Port, result[i-1].Port)
		}
	}
}

func TestSort_EmptySlice(t *testing.T) {
	s := portsorter.New(portsorter.ByPort, portsorter.Ascending)
	var entries []portsorter.Entry
	s.Sort(entries) // should not panic
}

package logger

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNew_Stdout(t *testing.T) {
	l, err := New("", "text")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer l.Close()

	if l.file != os.Stdout {
		t.Error("expected logger to write to stdout")
	}
}

func TestNew_InvalidFormat_DefaultsToText(t *testing.T) {
	l, err := New("", "xml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer l.Close()

	if l.format != "text" {
		t.Errorf("expected format 'text', got '%s'", l.format)
	}
}

func TestLog_TextFormat(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "portwatch-log-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	l, err := New(tmpFile.Name(), "text")
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer l.Close()

	e := Event{Timestamp: time.Now(), Port: 8080, Status: "open", Host: "localhost"}
	if err := l.Log(e); err != nil {
		t.Fatalf("Log returned error: %v", err)
	}

	data, _ := os.ReadFile(tmpFile.Name())
	contents := string(data)

	if !strings.Contains(contents, "port=8080") {
		t.Errorf("expected log to contain port=8080, got: %s", contents)
	}
	if !strings.Contains(contents, "status=open") {
		t.Errorf("expected log to contain status=open, got: %s", contents)
	}
}

func TestLog_JSONFormat(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "portwatch-log-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	l, err := New(tmpFile.Name(), "json")
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer l.Close()

	e := Event{Timestamp: time.Now(), Port: 443, Status: "closed", Host: "127.0.0.1"}
	if err := l.Log(e); err != nil {
		t.Fatalf("Log returned error: %v", err)
	}

	data, _ := os.ReadFile(tmpFile.Name())

	var decoded Event
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to parse JSON log line: %v", err)
	}
	if decoded.Port != 443 {
		t.Errorf("expected port 443, got %d", decoded.Port)
	}
	if decoded.Status != "closed" {
		t.Errorf("expected status 'closed', got '%s'", decoded.Status)
	}
}

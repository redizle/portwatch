package notify

import "testing"

func TestLevelFor_NewOpen(t *testing.T) {
	if got := LevelFor("open", ""); got != LevelInfo {
		t.Errorf("expected Info for new open port, got %s", got)
	}
}

func TestLevelFor_ReopenedPort(t *testing.T) {
	if got := LevelFor("open", "closed"); got != LevelWarn {
		t.Errorf("expected Warn for reopened port, got %s", got)
	}
}

func TestLevelFor_ClosedPort(t *testing.T) {
	if got := LevelFor("closed", "open"); got != LevelAlert {
		t.Errorf("expected Alert for closed port, got %s", got)
	}
}

func TestLevelFor_Default(t *testing.T) {
	if got := LevelFor("closed", "closed"); got != LevelInfo {
		t.Errorf("expected Info for no-change, got %s", got)
	}
}

func TestLevel_String(t *testing.T) {
	cases := []struct {
		level Level
		want  string
	}{
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelAlert, "ALERT"},
		{Level("unknown"), "UNKNOWN"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("Level(%q).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

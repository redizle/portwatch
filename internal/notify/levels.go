package notify

// LevelFor returns the appropriate notification level for a given port status transition.
// opened ports default to Warn; closed ports that were previously open use Alert;
// all other cases use Info.
func LevelFor(status, previous string) Level {
	switch {
	case status == "open" && previous == "":
		return LevelInfo
	case status == "open" && previous == "closed":
		return LevelWarn
	case status == "closed" && previous == "open":
		return LevelAlert
	default:
		return LevelInfo
	}
}

// String returns a human-readable label for a Level.
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelAlert:
		return "ALERT"
	default:
		return "UNKNOWN"
	}
}

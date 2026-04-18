package portaudit

import (
	"encoding/json"
	"io"
)

// Exporter writes audit entries to an io.Writer in a chosen format.
type Exporter struct {
	log *Log
}

// NewExporter wraps a Log for export.
func NewExporter(l *Log) *Exporter {
	return &Exporter{log: l}
}

// WriteJSON encodes all current entries as a JSON array.
func (e *Exporter) WriteJSON(w io.Writer) error {
	entries := e.log.All()
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

// WriteText writes entries in human-readable form, one per line.
func (e *Exporter) WriteText(w io.Writer) error {
	for _, entry := range e.log.All() {
		if _, err := io.WriteString(w, entry.String()+"\n"); err != nil {
			return err
		}
	}
	return nil
}

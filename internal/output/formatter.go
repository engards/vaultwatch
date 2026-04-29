package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/vaultwatch/internal/monitor"
)

// Format represents the output format type.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
)

// Formatter writes secret status output in a given format.
type Formatter struct {
	format Format
	w      io.Writer
}

// New creates a new Formatter. Defaults to stdout if w is nil.
func New(format Format, w io.Writer) *Formatter {
	if w == nil {
		w = os.Stdout
	}
	return &Formatter{format: format, w: w}
}

// Write renders a slice of SecretStatus to the configured output.
func (f *Formatter) Write(statuses []monitor.SecretStatus) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(statuses)
	default:
		return f.writeTable(statuses)
	}
}

func (f *Formatter) writeTable(statuses []monitor.SecretStatus) error {
	tw := tabwriter.NewWriter(f.w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "PATH\tSTATUS\tEXPIRES IN\tRENEWABLE")
	fmt.Fprintln(tw, "----\t------\t----------\t---------")
	for _, s := range statuses {
		expiry := formatDuration(s.TTL)
		renewable := "no"
		if s.Renewable {
			renewable = "yes"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", s.Path, s.Status, expiry, renewable)
	}
	return tw.Flush()
}

func (f *Formatter) writeJSON(statuses []monitor.SecretStatus) error {
	enc := json.NewEncoder(f.w)
	enc.SetIndent("", "  ")
	return enc.Encode(statuses)
}

func formatDuration(d time.Duration) string {
	if d <= 0 {
		return "expired"
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh%dm", h, m)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

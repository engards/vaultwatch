// Package diff compares two snapshots of secret statuses and reports changes.
package diff

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/your-org/vaultwatch/internal/monitor"
)

// ChangeType describes how a secret's status changed between snapshots.
type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeRemoved ChangeType = "removed"
	ChangeUpdated ChangeType = "updated"
)

// Change represents a single secret status change.
type Change struct {
	Path   string
	Type   ChangeType
	Before *monitor.SecretStatus
	After  *monitor.SecretStatus
}

// Differ computes differences between secret status snapshots.
type Differ struct {
	out io.Writer
}

// New returns a Differ that writes output to w.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Differ {
	if w == nil {
		w = os.Stdout
	}
	return &Differ{out: w}
}

// Compare returns the list of changes between the previous and current
// slices of SecretStatus, keyed by Path.
func (d *Differ) Compare(prev, curr []monitor.SecretStatus) []Change {
	prevMap := index(prev)
	currMap := index(curr)

	var changes []Change

	for path, after := range currMap {
		after := after
		if before, ok := prevMap[path]; !ok {
			changes = append(changes, Change{Path: path, Type: ChangeAdded, After: &after})
		} else if before.State != after.State || !before.ExpiresAt.Equal(after.ExpiresAt) {
			b := before
			changes = append(changes, Change{Path: path, Type: ChangeUpdated, Before: &b, After: &after})
		}
	}

	for path, before := range prevMap {
		if _, ok := currMap[path]; !ok {
			b := before
			changes = append(changes, Change{Path: path, Type: ChangeRemoved, Before: &b})
		}
	}

	return changes
}

// Print writes a human-readable summary of changes to the Differ's writer.
func (d *Differ) Print(changes []Change) {
	if len(changes) == 0 {
		fmt.Fprintln(d.out, "no changes detected")
		return
	}
	tw := tabwriter.NewWriter(d.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "CHANGE\tPATH\tSTATE")
	for _, c := range changes {
		state := ""
		if c.After != nil {
			state = string(c.After.State)
		} else if c.Before != nil {
			state = string(c.Before.State)
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\n", c.Type, c.Path, state)
	}
	tw.Flush()
}

func index(statuses []monitor.SecretStatus) map[string]monitor.SecretStatus {
	m := make(map[string]monitor.SecretStatus, len(statuses))
	for _, s := range statuses {
		m[s.Path] = s
	}
	return m
}

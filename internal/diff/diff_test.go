package diff_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/diff"
	"github.com/your-org/vaultwatch/internal/monitor"
)

var (
	now    = time.Now()
	later  = now.Add(24 * time.Hour)
	soonA  = monitor.SecretStatus{Path: "secret/a", State: monitor.StateOK, ExpiresAt: later}
	soonB  = monitor.SecretStatus{Path: "secret/b", State: monitor.StateWarning, ExpiresAt: now.Add(time.Hour)}
	soonC  = monitor.SecretStatus{Path: "secret/c", State: monitor.StateCritical, ExpiresAt: now.Add(10 * time.Minute)}
)

func TestCompare_NoChanges(t *testing.T) {
	d := diff.New(nil)
	changes := d.Compare([]monitor.SecretStatus{soonA}, []monitor.SecretStatus{soonA})
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}

func TestCompare_Added(t *testing.T) {
	d := diff.New(nil)
	changes := d.Compare(nil, []monitor.SecretStatus{soonA})
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != diff.ChangeAdded {
		t.Errorf("expected ChangeAdded, got %s", changes[0].Type)
	}
	if changes[0].Path != soonA.Path {
		t.Errorf("unexpected path %s", changes[0].Path)
	}
}

func TestCompare_Removed(t *testing.T) {
	d := diff.New(nil)
	changes := d.Compare([]monitor.SecretStatus{soonA}, nil)
	if len(changes) != 1 || changes[0].Type != diff.ChangeRemoved {
		t.Fatalf("expected 1 ChangeRemoved, got %v", changes)
	}
}

func TestCompare_Updated_StateChange(t *testing.T) {
	updated := soonA
	updated.State = monitor.StateWarning
	d := diff.New(nil)
	changes := d.Compare([]monitor.SecretStatus{soonA}, []monitor.SecretStatus{updated})
	if len(changes) != 1 || changes[0].Type != diff.ChangeUpdated {
		t.Fatalf("expected 1 ChangeUpdated, got %v", changes)
	}
	if changes[0].Before.State != monitor.StateOK {
		t.Errorf("expected before state OK, got %s", changes[0].Before.State)
	}
}

func TestCompare_MultipleChanges(t *testing.T) {
	prev := []monitor.SecretStatus{soonA, soonB}
	curr := []monitor.SecretStatus{soonB, soonC}
	d := diff.New(nil)
	changes := d.Compare(prev, curr)
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes (added+removed), got %d", len(changes))
	}
}

func TestPrint_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	d := diff.New(&buf)
	d.Print(nil)
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes' message, got: %s", buf.String())
	}
}

func TestPrint_WritesTable(t *testing.T) {
	var buf bytes.Buffer
	d := diff.New(&buf)
	changes := []diff.Change{
		{Path: "secret/a", Type: diff.ChangeAdded, After: &soonA},
	}
	d.Print(changes)
	out := buf.String()
	if !strings.Contains(out, "secret/a") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, "added") {
		t.Errorf("expected change type in output, got: %s", out)
	}
}

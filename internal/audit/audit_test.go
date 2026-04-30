package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/vaultwatch/internal/audit"
)

func TestLog_WritesJSONEvent(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	err := l.Log(audit.EventSecretChecked, "secret/myapp/db", "secret checked", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var e audit.Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &e); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}

	if e.Type != audit.EventSecretChecked {
		t.Errorf("expected type %q, got %q", audit.EventSecretChecked, e.Type)
	}
	if e.Path != "secret/myapp/db" {
		t.Errorf("expected path %q, got %q", "secret/myapp/db", e.Path)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLog_WithMeta(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	meta := map[string]string{"ttl": "3600", "lease_id": "abc-123"}
	_ = l.Log(audit.EventRenewalTriggered, "secret/myapp/token", "renewal triggered", meta)

	var e audit.Event
	_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &e)

	if e.Meta["lease_id"] != "abc-123" {
		t.Errorf("expected meta lease_id %q, got %q", "abc-123", e.Meta["lease_id"])
	}
}

func TestLog_MultipleEvents(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	_ = l.Log(audit.EventSecretChecked, "a", "msg1", nil)
	_ = l.Log(audit.EventSecretExpired, "b", "msg2", nil)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 log lines, got %d", len(lines))
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	l := audit.New(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

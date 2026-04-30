package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notify"
)

func TestSend_Success(t *testing.T) {
	var received notify.WebhookPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewWebhookNotifier(ts.URL, nil)
	payload := notify.WebhookPayload{
		Level:     "critical",
		Secret:    "secret/db/password",
		Message:   "expires soon",
		ExpiresIn: "2h",
		Timestamp: time.Now(),
	}

	if err := n.Send(payload); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Level != "critical" {
		t.Errorf("expected level critical, got %s", received.Level)
	}
	if received.Secret != "secret/db/password" {
		t.Errorf("expected secret path, got %s", received.Secret)
	}
}

func TestSend_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notify.NewWebhookNotifier(ts.URL, nil)
	err := n.Send(notify.WebhookPayload{Timestamp: time.Now()})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestSend_CustomHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if v := r.Header.Get("X-Api-Key"); v != "secret-key" {
			t.Errorf("expected X-Api-Key header, got %q", v)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	n := notify.NewWebhookNotifier(ts.URL, map[string]string{"X-Api-Key": "secret-key"})
	if err := n.Send(notify.WebhookPayload{Timestamp: time.Now()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewWebhookNotifier_Defaults(t *testing.T) {
	n := notify.NewWebhookNotifier("http://example.com", nil)
	if n.URL != "http://example.com" {
		t.Errorf("unexpected URL: %s", n.URL)
	}
	if n.Client == nil {
		t.Error("expected non-nil HTTP client")
	}
}

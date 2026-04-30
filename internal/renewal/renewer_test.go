package renewal_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/renewal"
)

func testLogger() *log.Logger {
	return log.New(os.Stdout, "", 0)
}

func TestNew_DefaultGracePct(t *testing.T) {
	r := renewal.New(nil, testLogger(), 0)
	// gracePct invalid → defaults to 0.2; ShouldRenew with 20% remaining should be true
	maxTTL := 100 * time.Second
	ttl := 20 * time.Second
	if !r.ShouldRenew(ttl, maxTTL) {
		t.Error("expected ShouldRenew=true at exactly grace boundary")
	}
}

func TestNew_InvalidGracePctClamped(t *testing.T) {
	r := renewal.New(nil, testLogger(), 1.5)
	// clamped to 0.2
	maxTTL := 100 * time.Second
	ttl := 19 * time.Second
	if !r.ShouldRenew(ttl, maxTTL) {
		t.Error("expected ShouldRenew=true when ttl < 20% of maxTTL")
	}
}

func TestShouldRenew_AboveThreshold(t *testing.T) {
	r := renewal.New(nil, testLogger(), 0.25)
	maxTTL := 200 * time.Second
	ttl := 100 * time.Second // 50% remaining — above 25% threshold
	if r.ShouldRenew(ttl, maxTTL) {
		t.Error("expected ShouldRenew=false when ttl is well above threshold")
	}
}

func TestShouldRenew_BelowThreshold(t *testing.T) {
	r := renewal.New(nil, testLogger(), 0.25)
	maxTTL := 200 * time.Second
	ttl := 40 * time.Second // 20% remaining — below 25% threshold
	if !r.ShouldRenew(ttl, maxTTL) {
		t.Error("expected ShouldRenew=true when ttl is below threshold")
	}
}

func TestShouldRenew_ZeroMaxTTL(t *testing.T) {
	r := renewal.New(nil, testLogger(), 0.2)
	if r.ShouldRenew(0, 0) {
		t.Error("expected ShouldRenew=false when maxTTL is zero (non-renewable)")
	}
}

// Package filter provides path-based and tag-based filtering of secret statuses.
package filter

import (
	"strings"

	"github.com/yourusername/vaultwatch/internal/monitor"
)

// Options holds the criteria used to filter secrets.
type Options struct {
	// PathPrefix restricts results to secrets whose path starts with this value.
	PathPrefix string
	// States restricts results to secrets in one of these states.
	// Valid values: "ok", "warning", "critical", "expired".
	States []string
}

// Filter applies Options to a slice of SecretStatus and returns matching entries.
type Filter struct {
	opts Options
}

// New creates a Filter with the given Options.
func New(opts Options) *Filter {
	return &Filter{opts: opts}
}

// Apply returns only those statuses that satisfy all configured criteria.
func (f *Filter) Apply(statuses []monitor.SecretStatus) []monitor.SecretStatus {
	var result []monitor.SecretStatus
	for _, s := range statuses {
		if !f.matchesPrefix(s) {
			continue
		}
		if !f.matchesState(s) {
			continue
		}
		result = append(result, s)
	}
	return result
}

func (f *Filter) matchesPrefix(s monitor.SecretStatus) bool {
	if f.opts.PathPrefix == "" {
		return true
	}
	return strings.HasPrefix(s.Path, f.opts.PathPrefix)
}

func (f *Filter) matchesState(s monitor.SecretStatus) bool {
	if len(f.opts.States) == 0 {
		return true
	}
	state := strings.ToLower(string(s.State))
	for _, allowed := range f.opts.States {
		if strings.ToLower(allowed) == state {
			return true
		}
	}
	return false
}

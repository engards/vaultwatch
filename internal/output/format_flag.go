package output

import (
	"fmt"
	"strings"
)

// FormatFlag implements the flag.Value interface for Format,
// allowing it to be used directly with the standard flag or cobra packages.
type FormatFlag struct {
	Value Format
}

// NewFormatFlag returns a FormatFlag with the default table format.
func NewFormatFlag() *FormatFlag {
	return &FormatFlag{Value: FormatTable}
}

// String returns the current format value as a string.
func (f *FormatFlag) String() string {
	return string(f.Value)
}

// Set parses and validates the provided format string.
func (f *FormatFlag) Set(s string) error {
	switch Format(strings.ToLower(s)) {
	case FormatTable:
		f.Value = FormatTable
	case FormatJSON:
		f.Value = FormatJSON
	default:
		return fmt.Errorf("unsupported format %q: must be one of [table, json]", s)
	}
	return nil
}

// Type returns the type name for use in help text (cobra compatibility).
func (f *FormatFlag) Type() string {
	return "format"
}

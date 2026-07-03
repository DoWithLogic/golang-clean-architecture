// Package logging provides a reusable Echo middleware for structured HTTP
// request/response logging via zerolog. It supports sensitive-field masking,
// path exclusion, and X-Request-ID propagation.
//
// # Quick Start
//
//	import (
//	    logging "github.com/golang-clean-architecture/pkg/logging"
//	)
//
//	e := echo.New()
//	e.Use(logging.Middleware(
//	    logging.WithLogger(&logger),
//	    logging.WithMaskedKeys("password", "token", "secret"),
//	))
package logging

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
)

// MaskedKey represents a JSON key whose value should be redacted in log output
// (for example: "password", "token", or "secret").
type MaskedKey string

// String returns the string representation of the MaskedKey.
func (m MaskedKey) String() string { return string(m) }

// IgnoredPatterns holds a collection of compiled regular expressions.
// Any request path matching at least one pattern is skipped during logging.
type IgnoredPatterns []*regexp.Regexp

// NewIgnoredPatterns compiles a slice of regex pattern strings into an
// IgnoredPatterns instance. It returns an error if any pattern is invalid.
//
// Example:
//
//	patterns, err := logging.NewIgnoredPatterns([]string{
//	    "^/metrics$",
//	    "^/ping$",
//	    "^/swagger/.*",
//	})
func NewIgnoredPatterns(patterns []string) (IgnoredPatterns, error) {
	compiled := make(IgnoredPatterns, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid ignored path pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return compiled, nil
}

// MatchesAny reports whether path matches at least one compiled pattern.
func (ip IgnoredPatterns) MatchesAny(path string) bool {
	for _, re := range ip {
		if re.MatchString(path) {
			return true
		}
	}
	return false
}

// loggingConfig configures the behaviour of the logging middleware.
type loggingConfig struct {
	// isObservabilityEnable toggles observability-specific log fields.
	isObservabilityEnable bool

	// maskedKeys lists keys whose values are replaced with "*****" in logs.
	// Matching is case-sensitive substring matching against the JSON key name.
	maskedKeys []MaskedKey

	// ignoredPatterns lists compiled regex patterns for paths that should be
	// excluded from logging entirely (for example: health-check or metrics endpoints).
	ignoredPatterns IgnoredPatterns

	// logger is the zerolog logger used by the middleware.
	logger *zerolog.Logger
}

func defaultConfig() loggingConfig {
	nop := zerolog.Nop()
	return loggingConfig{
		isObservabilityEnable: false,
		maskedKeys:            []MaskedKey{},
		ignoredPatterns:       IgnoredPatterns{},
		logger:                &nop,
	}
}

// IsPathIgnored reports whether the given request path matches any ignored pattern.
func (o *loggingConfig) IsPathIgnored(path string) bool {
	return o.ignoredPatterns.MatchesAny(path)
}

// IsMaskedKey reports whether the given key name should be redacted.
func (o *loggingConfig) IsMaskedKey(key string) bool {
	for _, mk := range o.maskedKeys {
		if strings.Contains(key, mk.String()) {
			return true
		}
	}
	return false
}

// MaskSensitiveFields recursively walks a JSON-decoded map and replaces
// values for any key matched by IsMaskedKey with "*****".
func (o *loggingConfig) MaskSensitiveFields(data map[string]any) map[string]any {
	for key, value := range data {
		switch {
		case o.IsMaskedKey(key):
			data[key] = "*****"
		case isNestedMap(value):
			data[key] = o.MaskSensitiveFields(value.(map[string]any))
		}
	}
	return data
}

func isNestedMap(v any) bool {
	_, ok := v.(map[string]any)
	return ok
}

// LoggingOption configures the logging middleware using the functional
// options pattern.
type LoggingOption func(*loggingConfig)

// WithLogger sets the zerolog logger used by the middleware.
//
// If l is nil, the default no-op logger is retained.
func WithLogger(l *zerolog.Logger) LoggingOption {
	return func(lc *loggingConfig) {
		if l != nil {
			lc.logger = l
		}
	}
}

// WithMaskedKeys sets the list of JSON keys whose values should be redacted
// in logs.
//
// Matching uses case-sensitive substring comparison against key names.
func WithMaskedKeys(mk ...string) LoggingOption {
	maskedKeys := make([]MaskedKey, 0, len(mk))
	for _, key := range mk {
		maskedKeys = append(maskedKeys, MaskedKey(key))
	}

	return func(lc *loggingConfig) { lc.maskedKeys = maskedKeys }
}

// WithIgnoredPatterns sets the request path patterns that should be excluded
// from logging.
func WithIgnoredPatterns(ips IgnoredPatterns) LoggingOption {
	return func(lc *loggingConfig) { lc.ignoredPatterns = ips }
}

// WithObservability enables or disables observability-specific log fields.
func WithObservability(enabled bool) LoggingOption {
	return func(lc *loggingConfig) { lc.isObservabilityEnable = enabled }
}

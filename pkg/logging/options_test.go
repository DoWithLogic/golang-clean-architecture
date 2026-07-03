package logging

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewIgnoredPatterns(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		patterns, err := NewIgnoredPatterns([]string{
			"^/ping$",
			"^/metrics$",
		})

		assert.NoError(t, err)
		assert.Len(t, patterns, 2)
		assert.True(t, patterns.MatchesAny("/ping"))
		assert.True(t, patterns.MatchesAny("/metrics"))
		assert.False(t, patterns.MatchesAny("/users"))
	})

	t.Run("invalid regex", func(t *testing.T) {
		patterns, err := NewIgnoredPatterns([]string{"["})

		assert.Error(t, err)
		assert.Nil(t, patterns)
	})
}

func TestIgnoredPatterns_MatchesAny(t *testing.T) {
	patterns, err := NewIgnoredPatterns([]string{
		"^/health$",
		"^/metrics$",
	})
	assert.NoError(t, err)

	assert.True(t, patterns.MatchesAny("/health"))
	assert.True(t, patterns.MatchesAny("/metrics"))
	assert.False(t, patterns.MatchesAny("/api/users"))
}

func TestWithLogger(t *testing.T) {
	logger := zerolog.Nop()
	cfg := defaultConfig()

	WithLogger(&logger)(&cfg)

	assert.Equal(t, &logger, cfg.logger)
}

func TestWithLogger_Nil(t *testing.T) {
	cfg := defaultConfig()
	original := cfg.logger

	WithLogger(nil)(&cfg)

	assert.Equal(t, original, cfg.logger)
}

func TestWithMaskedKeys(t *testing.T) {
	cfg := defaultConfig()

	WithMaskedKeys("password", "token", "secret")(&cfg)

	assert.Equal(t, []MaskedKey{
		"password",
		"token",
		"secret",
	}, cfg.maskedKeys)
}

func TestWithIgnoredPatterns(t *testing.T) {
	cfg := defaultConfig()

	patterns, err := NewIgnoredPatterns([]string{"^/ping$"})
	assert.NoError(t, err)

	WithIgnoredPatterns(patterns)(&cfg)

	assert.Equal(t, patterns, cfg.ignoredPatterns)
}

func TestWithObservability(t *testing.T) {
	cfg := defaultConfig()

	WithObservability(true)(&cfg)
	assert.True(t, cfg.isObservabilityEnable)

	WithObservability(false)(&cfg)
	assert.False(t, cfg.isObservabilityEnable)
}

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()

	assert.False(t, cfg.isObservabilityEnable)
	assert.Empty(t, cfg.maskedKeys)
	assert.Empty(t, cfg.ignoredPatterns)
	assert.NotNil(t, cfg.logger)
}

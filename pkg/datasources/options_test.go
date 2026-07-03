package datasources

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()

	if cfg.maxIdleConns != 2 {
		t.Errorf("expected maxIdleConns=10, got %d", cfg.maxIdleConns)
	}

	if cfg.maxOpenConns != 5 {
		t.Errorf("expected maxOpenConns=20, got %d", cfg.maxOpenConns)
	}

	if cfg.connMaxLifetime != time.Hour {
		t.Errorf("expected connMaxLifetime=1h, got %v", cfg.connMaxLifetime)
	}

	if cfg.connMaxIdleTime != 10*time.Minute {
		t.Errorf("expected connMaxIdleTime=10m, got %v", cfg.connMaxIdleTime)
	}
}

func TestOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		opts   []Option
		assert func(*testing.T, *config)
	}{
		{
			name: "WithMaxIdleConns",
			opts: []Option{WithMaxIdleConns(50)},
			assert: func(t *testing.T, c *config) {
				if c.maxIdleConns != 50 {
					t.Errorf("expected 50, got %d", c.maxIdleConns)
				}
			},
		},
		{
			name: "WithMaxOpenConns",
			opts: []Option{WithMaxOpenConns(200)},
			assert: func(t *testing.T, c *config) {
				if c.maxOpenConns != 200 {
					t.Errorf("expected 200, got %d", c.maxOpenConns)
				}
			},
		},
		{
			name: "WithConnMaxLifetime",
			opts: []Option{WithConnMaxLifetime(2 * time.Hour)},
			assert: func(t *testing.T, c *config) {
				if c.connMaxLifetime != 2*time.Hour {
					t.Errorf("expected 2h, got %v", c.connMaxLifetime)
				}
			},
		},
		{
			name: "WithConnMaxIdleTime",
			opts: []Option{WithConnMaxIdleTime(5 * time.Minute)},
			assert: func(t *testing.T, c *config) {
				if c.connMaxIdleTime != 5*time.Minute {
					t.Errorf("expected 5m, got %v", c.connMaxIdleTime)
				}
			},
		},
		{
			name: "WithMetrics",
			opts: []Option{WithMetrics(true)},
			assert: func(t *testing.T, c *config) {
				if c.enableMetrics != true {
					t.Errorf("expected true, got %v", c.enableMetrics)
				}
			},
		},
		{
			name: "WithTracing",
			opts: []Option{WithTracing(true)},
			assert: func(t *testing.T, c *config) {
				if c.enableTracing != true {
					t.Errorf("expected true, got %v", c.enableTracing)
				}
			},
		},
		{
			name: "multiple options combined",
			opts: []Option{
				WithMaxIdleConns(5),
				WithMaxOpenConns(10),
				WithConnMaxLifetime(30 * time.Minute),
				WithConnMaxIdleTime(1 * time.Minute),
				WithMetrics(true),
				WithTracing(true),
			},
			assert: func(t *testing.T, c *config) {
				if c.maxIdleConns != 5 {
					t.Errorf("expected 5, got %d", c.maxIdleConns)
				}
				if c.maxOpenConns != 10 {
					t.Errorf("expected 10, got %d", c.maxOpenConns)
				}
				if c.connMaxLifetime != 30*time.Minute {
					t.Errorf("expected 30m, got %v", c.connMaxLifetime)
				}
				if c.connMaxIdleTime != 1*time.Minute {
					t.Errorf("expected 1m, got %v", c.connMaxIdleTime)
				}
				if !c.enableMetrics {
					t.Errorf("expected metrics true")
				}
				if !c.enableTracing {
					t.Errorf("expected tracing true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := defaultConfig()

			for _, opt := range tt.opts {
				opt(cfg)
			}

			tt.assert(t, cfg)
		})
	}
}

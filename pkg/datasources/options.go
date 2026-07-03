package datasources

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

// DatabaseConfig holds the configuration for the database connection.
type DatabaseConfig struct {
	// Host is the address of the database server.
	Host string
	// Port is the port number of the database server.
	Port string
	// DBName is the name of the database to connect to.
	DBName string
	// UserName is the username for database authentication.
	UserName string
	// Password is the password for database authentication.
	Password string
	// Schema is the database schema (used for PostgreSQL).
	Schema string
	// Debug enables query debugging when set to true.
	Debug bool
}

func (d DatabaseConfig) toMYSQLConfig() mysql.Config {
	return mysql.Config{
		User:                 d.UserName,
		Passwd:               d.Password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", d.Host, d.Port),
		DBName:               d.DBName,
		ParseTime:            true,
		Loc:                  time.UTC,
		AllowNativePasswords: true,
	}
}

// config holds internal configuration for database connection setup.
//
// It is built using the functional options pattern and should not be
// instantiated directly outside this package.
type config struct {
	driver string
	dsn    string
	debug  bool
	schema string

	maxIdleConns    int
	maxOpenConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration

	enableMetrics bool
	enableTracing bool
	isProduction  bool
}

func NewGORMOptions(isObservabilityEnable, isProduction bool) []Option {
	return []Option{
		WithMetrics(isObservabilityEnable),
		WithTracing(isObservabilityEnable),
		IsProduction(isProduction),
	}
}

// defaultConfig returns a config initialized with safe default values.
//
// Default values are optimized for general-purpose production usage:
//   - moderate connection pool sizes
//   - reasonable connection lifetime limits
//
// These defaults can be overridden using Option functions.
func defaultConfig() *config {
	return &config{
		maxIdleConns:    2,
		maxOpenConns:    5,
		connMaxLifetime: 1 * time.Hour,
		connMaxIdleTime: 10 * time.Minute,
	}
}

func (c *config) nonProductionConnectionPool() {
	c.maxIdleConns = 1
	c.maxOpenConns = 2
}

// Option represents a functional option used to modify config.
//
// Options are applied in order during database initialization.
// This allows flexible and composable configuration without
// breaking constructor signatures.
type Option func(*config)

// WithMaxIdleConns sets the maximum number of idle connections
// in the database connection pool.
//
// A higher value improves reuse but consumes more resources.
func WithMaxIdleConns(n int) Option {
	return func(c *config) {
		c.maxIdleConns = n
	}
}

// WithMaxOpenConns sets the maximum number of open connections
// to the database.
//
// This controls the total concurrency allowed at the DB level.
func WithMaxOpenConns(n int) Option {
	return func(c *config) {
		c.maxOpenConns = n
	}
}

// WithConnMaxLifetime sets the maximum amount of time a connection
// may be reused.
//
// Connections older than this will be closed and replaced.
func WithConnMaxLifetime(d time.Duration) Option {
	return func(c *config) {
		c.connMaxLifetime = d
	}
}

// WithConnMaxIdleTime sets the maximum amount of time a connection
// may remain idle in the pool before being closed.
//
// This helps reduce resource usage in low traffic scenarios.
func WithConnMaxIdleTime(d time.Duration) Option {
	return func(c *config) {
		c.connMaxIdleTime = d
	}
}

// WithMetrics enables or disables Prometheus database metrics collection.
//
// When enabled, database connection statistics will be exported
// via the Prometheus client.
func WithMetrics(enable bool) Option {
	return func(c *config) {
		c.enableMetrics = enable
	}
}

// WithTracing enables or disables OpenTelemetry tracing support.
//
// When enabled, database queries will be traced automatically
// using the configured tracing plugin.
func WithTracing(enable bool) Option {
	return func(c *config) {
		c.enableTracing = enable
	}
}

func IsProduction(isProd bool) Option {
	return func(c *config) { c.isProduction = isProd }
}

package datasources

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

func NewPostgresDB(ctx context.Context, cfg DatabaseConfig, opts ...Option) (*gorm.DB, error) {
	c := defaultConfig()
	c.debug = cfg.Debug
	c.schema = cfg.Schema

	c.dsn = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s search_path=%s sslmode=disable TimeZone=UTC",
		cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Schema,
	)

	for _, opt := range opts {
		opt(c)
	}

	if !c.isProduction {
		c.nonProductionConnectionPool()
	}

	gormCfg := &gorm.Config{TranslateError: true}
	if !c.debug {
		gormCfg.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(c.dsn), gormCfg)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(c.maxIdleConns)
	sqlDB.SetMaxOpenConns(c.maxOpenConns)
	sqlDB.SetConnMaxLifetime(c.connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(c.connMaxIdleTime)

	if c.enableMetrics {
		prometheus.MustRegister(collectors.NewDBStatsCollector(sqlDB, cfg.DBName))
	}

	if c.enableTracing {
		_ = db.Use(tracing.NewPlugin())
	}

	return db, nil
}

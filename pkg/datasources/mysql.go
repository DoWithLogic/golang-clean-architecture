package datasources

import (
	"fmt"

	"github.com/DoWithLogic/golang-clean-architecture/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	otps := &gorm.Config{SkipDefaultTransaction: true}

	if !cfg.Debug {
		otps = &gorm.Config{SkipDefaultTransaction: true, Logger: logger.Default.LogMode(logger.Silent)}
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Asia%%2FJakarta", cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	fmt.Print(dsn)
	db, err := gorm.Open(mysql.Open(dsn), otps)
	if err != nil {
		panic(err)
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}

	if !cfg.Debug {
		db.Logger.LogMode(logger.Silent)
	}

	return db, nil
}

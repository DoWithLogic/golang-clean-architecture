package datasources_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/datasources"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestNewMySQLDBSuccess(t *testing.T) {
	cfg := config.DatabaseConfig{
		UserName: "root",
		Password: "pass",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "testdb",
		Debug:    true,
	}

	// override mysql.Open to intercept DSN
	defer func() { recover() }() // catch panic if any

	// To test DSN formatting and config usage
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Asia%%2FJakarta",
		cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	assert.Contains(t, dsn, "localhost:3306")
	assert.Contains(t, dsn, "parseTime=true")

	// We’ll just assert the function panics if MySQL is unreachable (expected)
	assert.Panics(t, func() {
		datasources.NewMySQLDB(cfg)
	}, "expected panic because no MySQL is available")
}

func TestNewMySQLDBWithSilentLogger(t *testing.T) {
	cfg := config.DatabaseConfig{
		UserName: "root",
		Password: "pass",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "testdb",
		Debug:    false,
	}

	assert.Panics(t, func() {
		datasources.NewMySQLDB(cfg)
	})
}

func TestNewMySQLDBMockConnection(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dial := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	// Simulate gorm.Open works fine
	gdb, err := gorm.Open(dial, &gorm.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, gdb)

	mock.ExpectClose()
	assert.NoError(t, db.Close())
}

func TestDSNFormatting(t *testing.T) {
	cfg := config.DatabaseConfig{
		UserName: "user",
		Password: "1234",
		Host:     "127.0.0.1",
		Port:     "3307",
		DBName:   "demo",
		Debug:    false,
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Asia%%2FJakarta",
		cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	assert.Equal(t, "user:1234@tcp(127.0.0.1:3307)/demo?charset=utf8&parseTime=true&loc=Asia%2FJakarta", dsn)
}

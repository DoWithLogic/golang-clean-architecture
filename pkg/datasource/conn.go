package datasource

import (
	"fmt"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/jmoiron/sqlx"
)

func NewDatabase(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	// Create the connection string.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	// Connect to the MySQL database.
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Test the database connection.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

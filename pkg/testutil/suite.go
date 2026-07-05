package testutil

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

const (
	seedVersionTable      = "_goose_data_version"
	migrationVersionTable = "_goose_db_version"
)

type DBSuite struct {
	suite.Suite

	Ctx context.Context
	DB  *gorm.DB

	DatabaseName  string
	MigrationsDir string
	SeedsDir      string

	Dialect goose.Dialect
}

func (s *DBSuite) setupGoose(db *sql.DB) {
	if s.MigrationsDir != "" {
		goose.SetDialect(string(s.Dialect))
		goose.SetTableName(migrationVersionTable)

		s.Require().NoError(
			goose.UpContext(s.Ctx, db, s.MigrationsDir, goose.WithAllowMissing()),
		)
	}

	if s.SeedsDir != "" {
		goose.SetDialect(string(s.Dialect))
		goose.SetTableName(seedVersionTable)

		s.Require().NoError(
			goose.UpContext(s.Ctx, db, s.SeedsDir, goose.WithAllowMissing()),
		)
	}
}

func (s *DBSuite) teardownGoose(db *sql.DB) {
	if s.SeedsDir != "" {
		goose.SetDialect(string(s.Dialect))
		goose.SetTableName(seedVersionTable)

		_ = goose.DownToContext(s.Ctx, db, s.SeedsDir, 0, goose.WithAllowMissing())
	}

	if s.MigrationsDir != "" {
		goose.SetDialect(string(s.Dialect))
		goose.SetTableName(migrationVersionTable)

		_ = goose.DownToContext(s.Ctx, db, s.MigrationsDir, 0, goose.WithAllowMissing())
	}
}

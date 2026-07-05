package testutil

import (
	"context"
	"os"

	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/log"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultPGSQLImage = "postgres:15-alpine"

type PostgresSuite struct {
	DBSuite

	Container *tcpostgres.PostgresContainer
}

func (s *PostgresSuite) Setup() {
	s.Ctx = context.Background()

	testcontainers.WithLogger(log.Default())

	s.Dialect = goose.DialectPostgres

	image := os.Getenv("POSTGRES_CONTAINER_IMAGE")
	if image == "" {
		image = defaultPGSQLImage
	}

	var err error

	s.Container, err = tcpostgres.Run(s.Ctx, image, tcpostgres.WithDatabase(s.DatabaseName), tcpostgres.WithUsername("root"), tcpostgres.WithPassword("password"), tcpostgres.BasicWaitStrategies())
	s.Require().NoError(err)

	dsn, err := s.Container.ConnectionString(s.Ctx, "sslmode=disable", "search_path=public", "TimeZone=UTC")
	s.Require().NoError(err)

	s.DB, err = gorm.Open(gormpostgres.Open(dsn), &gorm.Config{TranslateError: true})
	s.Require().NoError(err)

	db, err := s.DB.DB()
	s.Require().NoError(err)

	s.setupGoose(db)
}

func (s *PostgresSuite) TearDown() {
	db, err := s.DB.DB()
	s.Require().NoError(err)

	s.teardownGoose(db)

	testcontainers.CleanupContainer(s.T(), s.Container.Container)
}

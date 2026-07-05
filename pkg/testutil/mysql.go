package testutil

import (
	"context"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const defaultMYSQLImage = "mysql:8.0.36"

type MYSQLSuite struct {
	DBSuite
	Container *tcmysql.MySQLContainer
}

func (s *MYSQLSuite) Setup() {
	s.Ctx = context.Background()
	s.Dialect = goose.DialectMySQL

	image := os.Getenv("MYSQL_CONTAINER_IMAGE")
	if image == "" {
		image = defaultMYSQLImage
	}

	var err error

	s.Container, err = tcmysql.Run(
		s.Ctx,
		image,
		tcmysql.WithDatabase(s.DatabaseName),
		tcmysql.WithUsername("root"),
		tcmysql.WithPassword("password"),
	)
	s.Require().NoError(err)

	dsn, err := s.Container.ConnectionString(s.Ctx, "charset=utf8", "parseTime=true", "loc=UTC", "multiStatements=true")
	s.Require().NoError(err)

	s.DB, err = gorm.Open(gormmysql.Open(dsn))
	s.Require().NoError(err)

	db, err := s.DB.DB()
	s.Require().NoError(err)

	s.setupGoose(db)
}

func (s *MYSQLSuite) TearDown() {
	db, err := s.DB.DB()
	s.Require().NoError(err)

	s.teardownGoose(db)

	testcontainers.CleanupContainer(s.T(), s.Container.Container)
}

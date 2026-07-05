package tests

import (
	"testing"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

type MYSQLRepositoryTestSuite struct {
	testutil.MYSQLSuite
}

func TestMYSQLRepositoryTestSuite(t *testing.T) { suite.Run(t, new(MYSQLRepositoryTestSuite)) }

func (s *MYSQLRepositoryTestSuite) SetupSuite() {
	s.DatabaseName = "golang_clean_architecture"
	s.MigrationsDir = "../database/mysql/migration"
	s.Setup()
}

func (s *MYSQLRepositoryTestSuite) TearDownSuite() { s.TearDown() }

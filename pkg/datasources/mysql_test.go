package datasources_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/datasources"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func startMySQL(t *testing.T) (*mysql.MySQLContainer, string) {
	ctx := context.Background()

	img := os.Getenv("MYSQL_CONTAINER_IMAGE")
	if img == "" {
		img = "mysql:8.0.36"
	}

	container, err := mysql.Run(
		ctx,
		img,
		mysql.WithDatabase("golang_clean_architecture"),
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
	)

	require.NoError(t, err)

	dsn, err := container.ConnectionString(
		ctx,
		"charset=utf8",
		"multiStatements=true",
		"loc=UTC",
		"parseTime=true",
	)
	require.NoError(t, err)

	return container, dsn
}

func TestNewMySQLDB_ConnectAndPing(t *testing.T) {
	t.Parallel()

	container, dsn := startMySQL(t)
	defer func() {
		_ = container.Terminate(context.Background())
	}()

	host, port, dbName := parseMySQLDSN(t, dsn)

	cfg := datasources.DatabaseConfig{
		Host:     host,
		Port:     port,
		DBName:   dbName,
		UserName: "root",
		Password: "password",
		Debug:    true,
	}

	db, err := datasources.NewMySQLDB(
		context.Background(),
		cfg,
		datasources.WithMaxOpenConns(10),
		datasources.WithMaxIdleConns(5),
		datasources.WithMetrics(false),
		datasources.WithTracing(false),
	)

	require.NoError(t, err)
	require.NotNil(t, db)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	require.NoError(t, sqlDB.Ping())

	stats := sqlDB.Stats()
	require.Equal(t, 2, stats.MaxOpenConnections)
}

func parseMySQLDSN(t *testing.T, dsn string) (host, port, db string) {
	require.Contains(t, dsn, "@tcp(")

	parts := strings.Split(dsn, "@tcp(")
	require.Len(t, parts, 2)

	hostPortDB := strings.Split(parts[1], ")")
	require.Len(t, hostPortDB, 2)

	hostPort := strings.Split(hostPortDB[0], ":")

	dbPart := strings.Split(hostPortDB[1], "/")
	dbName := strings.Split(dbPart[1], "?")[0]

	return hostPort[0], hostPort[1], dbName
}

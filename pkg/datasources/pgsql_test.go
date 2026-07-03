package datasources_test

import (
	"context"
	"strings"
	"testing"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/datasources"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	defaultImg = "postgres:15-alpine"
)

func startPostgres(t *testing.T) (*postgres.PostgresContainer, string) {
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		defaultImg,
		postgres.WithDatabase("golang_clean_architecture"),
		postgres.WithUsername("root"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
	)

	require.NoError(t, err)

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	return container, dsn
}

func TestNewPostgresDB_ConnectsSuccessfully(t *testing.T) {
	t.Parallel()

	container, dsn := startPostgres(t)
	defer func() {
		_ = container.Terminate(context.Background())
	}()

	host, port := parseDSN(t, dsn)

	cfg := datasources.DatabaseConfig{
		Host:     host,
		Port:     port,
		DBName:   "golang_clean_architecture",
		UserName: "root",
		Password: "password",
		Schema:   "public",
		Debug:    true,
	}

	db, err := datasources.NewPostgresDB(
		context.Background(),
		cfg,
		datasources.WithMaxOpenConns(5),
		datasources.WithMaxIdleConns(2),
		datasources.WithMetrics(false),
		datasources.WithTracing(false),
	)

	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	require.NoError(t, sqlDB.Ping())
}

func parseDSN(t *testing.T, dsn string) (host, port string) {
	parts := strings.Split(dsn, "@")
	require.Len(t, parts, 2)

	hostPort := strings.Split(strings.Split(parts[1], "/")[0], ":")
	require.Len(t, hostPort, 2)

	return hostPort[0], hostPort[1]
}

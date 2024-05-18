package dataprovider

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func preparePostgresContainer(t *testing.T) (*ConfigModule, testcontainers.Container) {
	ctx := context.Background()

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.3"),
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		t.Fatalf("Could not start mysql container: %s", err)
	}

	// Get the container's host and port
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Could not get mysql container host: %s", err)
	}

	netPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Could not get mysql container port: %s", err)
	}

	return &ConfigModule{
		Driver:   PostgreSQLDatabaseProviderName,
		Username: "test",
		Password: "test",
		Name:     "test",
		Host:     host,
		Port:     netPort.Int(),
	}, postgresContainer
}

func TestNewPostgresDataProvider(t *testing.T) {
	testCfg, container := preparePostgresContainer(t)
	defer container.Terminate(context.Background())

	provider, err := NewDataProvider(testCfg)
	assert.NoError(t, err)

	providerStatus := provider.GetProviderStatus()
	assert.Equal(t, PostgreSQLDatabaseProviderName, providerStatus.Driver)

	query, err := GetQueryFromFile("testdata/postgres/create_user_table.sql")
	assert.NoError(t, err)

	if err = provider.InitializeDatabase(query); err != nil {
		panic(err)
	}

	conn := provider.GetConnection()

	query, err = GetQueryFromFile("testdata/postgres/insert_user.sql")
	assert.NoError(t, err)

	tx := conn.MustBegin()
	tx.MustExec(query, "83.121.11.105", "New York")
	tx.MustExec(query, "76.71.94.89", "Los Angeles")
	tx.MustExec(query, "204.195.163.16", "Chicago")

	err = tx.Commit()
	assert.NoError(t, err)

	var users []struct {
		ID        int64  `json:"id" db:"id"`
		IpAddress string `json:"ip_address" db:"ip_address"`
		City      string `json:"city" db:"city"`
	}

	query, err = GetQueryFromFile("testdata/postgres/select_users.sql")
	assert.NoError(t, err)

	rows, err := conn.Queryx(query)
	assert.NoError(t, err)

	for rows.Next() {
		var user struct {
			ID        int64  `json:"id" db:"id"`
			IpAddress string `json:"ip_address" db:"ip_address"`
			City      string `json:"city" db:"city"`
		}
		if err = rows.StructScan(&user); err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	assert.Equalf(t, 3, len(users), "Expected 3 users, got %d", len(users))
}

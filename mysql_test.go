package dataprovider

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"testing"
)

func prepareMysqlContainer(t *testing.T) (*ConfigModule, testcontainers.Container) {
	ctx := context.Background()

	mysqContainer, err := mysql.RunContainer(ctx, testcontainers.WithImage("mysql:8.4"),
		mysql.WithDatabase("test"),
		mysql.WithUsername("test"),
		mysql.WithPassword("test"),
	)

	if err != nil {
		t.Fatalf("Could not start mysql container: %s", err)
	}

	// Get the container's host and port
	host, err := mysqContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Could not get mysql container host: %s", err)
	}

	netPort, err := mysqContainer.MappedPort(ctx, "3306")
	if err != nil {
		t.Fatalf("Could not get mysql container port: %s", err)
	}

	cfg := NewConfigModule().
		WithDriver(MySQLDatabaseProviderName).
		WithUsername("test").
		WithPassword("test").
		WithName("test").
		WithHost(host).
		WithPort(netPort.Int()).
		Build()

	return cfg, mysqContainer
}

func TestNewMySQLDataProvider(t *testing.T) {
	testCfg, container := prepareMysqlContainer(t)
	defer func(container testcontainers.Container, ctx context.Context) {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}(container, context.Background())

	provider, err := NewDataProvider(testCfg)
	assert.NoError(t, err)

	providerStatus := provider.GetProviderStatus()
	assert.Equal(t, MySQLDatabaseProviderName, providerStatus.Driver)

	query, err := GetQueryFromFile("testdata/mysql/create_user_table.sql")
	assert.NoError(t, err)

	if err = provider.InitializeDatabase(query); err != nil {
		panic(err)
	}

	conn := provider.GetConnection()

	query, err = GetQueryFromFile("testdata/mysql/insert_user.sql")
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

	query, err = GetQueryFromFile("testdata/mysql/select_users.sql")
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

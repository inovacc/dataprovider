package dataprovider

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func prepareOracleContainer(t *testing.T) (*ConfigModule, testcontainers.Container) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "gvenzl/oracle-free",
		ExposedPorts: []string{"1521/tcp"},
		Env:          map[string]string{"ORACLE_PASSWORD": "oracle"},
		WaitingFor:   wait.ForLog("Database ready to use. Enjoy! ;)").WithStartupTimeout(10 * time.Minute),
	}

	oracleContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Fatalf("Could not start oracle container: %s", err)
	}

	// Get the container's host and port
	host, err := oracleContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Could not get oracle container host: %s", err)
	}

	netPort, err := oracleContainer.MappedPort(ctx, "1521")
	if err != nil {
		t.Fatalf("Could not get oracle container port: %s", err)
	}

	return &ConfigModule{
		Driver:   OracleDatabaseProviderName,
		Username: "system",
		Password: "oracle",
		Name:     "xe",
		Host:     host,
		Port:     netPort.Int(),
	}, oracleContainer
}

func TestNewOracleDataProvider(t *testing.T) {
	testCfg, container := prepareOracleContainer(t)
	defer container.Terminate(context.Background())

	provider, err := NewDataProvider(testCfg)
	assert.NoError(t, err)

	providerStatus := provider.GetProviderStatus()
	assert.Equal(t, OracleDatabaseProviderName, providerStatus.Driver)

	query, err := GetQueryFromFile("testdata/oracle/create_user_table.sql")
	assert.NoError(t, err)

	if err = provider.InitializeDatabase(query); err != nil {
		panic(err)
	}

	conn := provider.GetConnection()

	query, err = GetQueryFromFile("testdata/oracle/insert_user.sql")
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

	query, err = GetQueryFromFile("testdata/oracle/select_users.sql")
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

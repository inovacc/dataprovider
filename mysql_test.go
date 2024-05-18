package dataprovider

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"log"
	"os"
	"testing"
)

var testCfg *ConfigModule

func TestMain(m *testing.M) {
	ctx := context.Background()

	mysqContainer, err := mysql.RunContainer(ctx, testcontainers.WithImage("mysql:8.4"),
		mysql.WithDatabase("test"),
		mysql.WithUsername("test"),
		mysql.WithPassword("test"),
	)

	if err != nil {
		log.Fatalf("Could not start mysql container: %s", err)
	}
	// Clean up the container
	defer func() {
		if err = mysqContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// Get the container's host and port
	host, err := mysqContainer.Host(ctx)
	if err != nil {
		log.Fatalf("Could not get mysql container host: %s", err)
	}

	netPort, err := mysqContainer.MappedPort(ctx, "3306")
	if err != nil {
		log.Fatalf("Could not get mysql container port: %s", err)
	}

	testCfg = &ConfigModule{
		Driver:   MySQLDatabaseProviderName,
		Username: "test",
		Password: "test",
		Name:     "test",
		Host:     host,
		Port:     netPort.Int(),
	}

	code := m.Run()
	os.Exit(code)
}

func TestNewDataProvider(t *testing.T) {
	provider, err := NewDataProvider(testCfg)
	assert.NoError(t, err)

	providerStatus := provider.GetProviderStatus()
	assert.Equal(t, MySQLDatabaseProviderName, providerStatus.Driver)

	query, err := GetQueryFromFile("testdata/create_user_table.sql")
	assert.NoError(t, err)

	if err = provider.InitializeDatabase(query); err != nil {
		panic(err)
	}

	conn := provider.GetConnection()

	query, err = GetQueryFromFile("testdata/insert_user.sql")
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

	query, err = GetQueryFromFile("testdata/select_users.sql")
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

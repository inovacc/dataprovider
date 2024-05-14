//go:build windows

package dataprovider

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMemoryDataProvider(t *testing.T) {
	cfg := &ConfigModule{
		Driver: MemoryDataProviderName,
	}

	provider, err := NewDataProvider(context.Background(), cfg)
	assert.NoError(t, err)

	providerStatus := provider.GetProviderStatus()
	assert.Equal(t, MemoryDataProviderName, providerStatus.Driver)
}

func TestNewSQLiteDataProvider(t *testing.T) {
	cfg := &ConfigModule{
		Driver:           SQLiteDataProviderName,
		ConnectionString: "file:test.sqlite3?cache=shared",
	}

	provider, err := NewDataProvider(context.Background(), cfg)
	assert.NoError(t, err)

	query := "CREATE TABLE IF NOT EXISTS test_table (id INTEGER PRIMARY KEY, name TEXT);"

	err = provider.InitializeDatabase(query)
	assert.NoError(t, err)

	conn := provider.GetConnection()
	assert.NotNil(t, conn)

	query = "INSERT INTO test_table (name) VALUES ('test');"

	_, err = conn.Exec(query)
	assert.NoError(t, err)

	user := struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}{}

	err = conn.Get(&user, "SELECT * FROM test_table")
	assert.NoError(t, err)

	assert.Equal(t, 1, user.ID)

	err = provider.Disconnect()
	assert.NoError(t, err)
}

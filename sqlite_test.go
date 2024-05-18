package dataprovider

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewSQLiteDataProvider(t *testing.T) {
	cfg := &ConfigModule{
		Driver:           SQLiteDataProviderName,
		ConnectionString: "file:test.sqlite3?cache=shared",
	}

	provider, err := NewDataProvider(cfg)
	assert.NoError(t, err)

	providerStatus := provider.GetProviderStatus()
	assert.Equal(t, SQLiteDataProviderName, providerStatus.Driver)

	query, err := GetQueryFromFile("testdata/sqlite/create_user_table.sql")
	assert.NoError(t, err)

	if err = provider.InitializeDatabase(query); err != nil {
		panic(err)
	}

	conn := provider.GetConnection()

	query, err = GetQueryFromFile("testdata/sqlite/insert_user.sql")
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

	query, err = GetQueryFromFile("testdata/sqlite/select_users.sql")
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

	err = provider.Disconnect()
	assert.NoError(t, err)

	err = os.Remove("test.sqlite3")
	assert.NoError(t, err)
}

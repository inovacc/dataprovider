package dataprovider

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMemoryProvider(t *testing.T) {
	driverName := SQLiteDataProviderName
	opts := NewOptions(WithDriver(driverName), WithConnectionString("file:test.sqlite3?cache=shared"))
	var provider = Must(NewDataProvider(opts))

	if providerStatus := provider.GetProviderStatus(); providerStatus.Driver != driverName {
		t.Errorf("Expected %s, got %s", driverName, providerStatus.Driver)
	}

	conn := provider.GetConnection()

	query, err := GetQueryFromFile("internal/testdata/sqlite/create_user_table.sql")
	if err != nil {
		t.Errorf("Error %s", err)
	}

	if err = provider.InitializeDatabase(query); err != nil {
		panic(err)
	}

	query, err = GetQueryFromFile("internal/testdata/sqlite/insert_user.sql")
	if err != nil {
		t.Errorf("Error %s", err)
	}

	tx := conn.MustBegin()
	tx.MustExec(query, "83.121.11.105", "New York")
	tx.MustExec(query, "76.71.94.89", "Los Angeles")
	tx.MustExec(query, "204.195.163.16", "Chicago")

	err = tx.Commit()
	assert.NoError(t, err)
}

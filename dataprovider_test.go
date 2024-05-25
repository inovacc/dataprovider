package dataprovider

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPostgresDataProvider(t *testing.T) {
	opts := NewOptions()
	var postgresProvider = Must(NewDataProvider(opts))

	providerStatus := postgresProvider.GetProviderStatus()
	if providerStatus.Driver != MemoryDataProviderName {
		t.Errorf("Expected %s, got %s", MemoryDataProviderName, providerStatus.Driver)
	}

	conn := postgresProvider.GetConnection()

	query, err := GetQueryFromFile("testdata/sqlite/create_user_table.sql")
	if err != nil {
		t.Errorf("Error %s", err)
	}

	if err = postgresProvider.InitializeDatabase(query); err != nil {
		panic(err)
	}

	query, err = GetQueryFromFile("testdata/sqlite/insert_user.sql")
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

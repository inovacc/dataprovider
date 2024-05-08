package dataprovider

import (
	"context"
	"github.com/jmoiron/sqlx"
)

// MemoryProvider defines the auth provider for in-memory database
type MemoryProvider struct {
	dbHandle *sqlx.DB
}

func (m *MemoryProvider) Connect() error {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryProvider) Disconnect() error {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryProvider) GetConnection() *sqlx.DB {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryProvider) CheckAvailability() error {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryProvider) ReconnectDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryProvider) InitializeDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryProvider) migrateDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryProvider) RevertDatabase(targetVersion int) error {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryProvider) resetDatabase() error {
	//TODO implement me
	panic("implement me")
}

func newMemoryProvider(ctx context.Context) error {
	dbHandle, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}

	if err = dbHandle.PingContext(ctx); err != nil {
		return err
	}

	provider = &MemoryProvider{dbHandle: dbHandle}

	return nil
}

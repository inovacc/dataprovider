package dataprovider

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// MemoryProvider defines the auth provider for in-memory database
type MemoryProvider struct {
	dbHandle *sqlx.DB
}

func (m *MemoryProvider) MigrateDatabase() error {
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
	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()

	return m.dbHandle.PingContext(ctx)
}

func (m *MemoryProvider) ReconnectDatabase() error {
	return m.CheckAvailability()
}

func (m *MemoryProvider) InitializeDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryProvider) RevertDatabase(targetVersion int) error {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryProvider) ResetDatabase() error {
	//TODO implement me
	panic("implement me")
}

func newMemoryProvider(ctx context.Context, cfg *ConfigModule) (*Wrapper, error) {
	dbHandle, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	if err = dbHandle.PingContext(ctx); err != nil {
		return nil, err
	}

	return &Wrapper{
		Driver:   cfg.Driver,
		Version:  1,
		Provider: &MemoryProvider{dbHandle: dbHandle},
	}, nil
}

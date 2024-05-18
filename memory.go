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

// GetProviderStatus returns the status of the provider
func (m *MemoryProvider) GetProviderStatus() ProviderStatus {
	status := ProviderStatus{
		Driver:   driverName,
		IsActive: true,
	}

	if err := m.CheckAvailability(); err != nil {
		status.IsActive = false
		status.Error = err
	}

	return status
}

// MigrateDatabase migrates the database to the latest version
func (m *MemoryProvider) MigrateDatabase() error {
	//TODO implement me
	panic("implement me")
}

// Disconnect disconnects from the data provider
func (m *MemoryProvider) Disconnect() error {
	return m.dbHandle.Close()
}

// GetConnection returns the connection to the data provider
func (m *MemoryProvider) GetConnection() *sqlx.DB {
	return m.dbHandle
}

// CheckAvailability checks if the data provider is available
func (m *MemoryProvider) CheckAvailability() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()

	return m.dbHandle.PingContext(ctx)
}

// ReconnectDatabase reconnects to the database
func (m *MemoryProvider) ReconnectDatabase() error {
	return m.CheckAvailability()
}

// InitializeDatabase initializes the database
func (m *MemoryProvider) InitializeDatabase(schema string) error {
	_, err := m.dbHandle.Exec(schema)
	return err
}

// MigrateDatabase migrates the database to the latest version
func (m *MemoryProvider) RevertDatabase(targetVersion int) error {
	//TODO implement me
	panic("implement me")
}

// ResetDatabase resets the database
func (m *MemoryProvider) ResetDatabase() error {
	//TODO implement me
	panic("implement me")
}

// newMemoryProvider creates a new memory provider
func newMemoryProvider(ctx context.Context, cfg *ConfigModule) (*MemoryProvider, error) {
	dbHandle, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	if err = dbHandle.PingContext(ctx); err != nil {
		return nil, err
	}

	return &MemoryProvider{dbHandle: dbHandle}, nil
}

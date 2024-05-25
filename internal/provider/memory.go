package provider

import (
	"context"
	"github.com/dyammarcano/dataprovider/internal/migration"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// MemoryProvider defines the auth provider for in-memory database
type MemoryProvider struct {
	dbHandle *sqlx.DB
}

// GetProviderStatus returns the status of the provider
func (m *MemoryProvider) GetProviderStatus() Status {
	status := Status{
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
func (m *MemoryProvider) MigrateDatabase() migration.MigrationProvider {
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

// RevertDatabase migrates the database to the latest version
func (m *MemoryProvider) RevertDatabase(targetVersion int) error {
	//TODO implement me
	panic("implement me")
}

// ResetDatabase resets the database
func (m *MemoryProvider) ResetDatabase() error {
	//TODO implement me
	panic("implement me")
}

// NewMemoryProvider creates a new memory provider
func NewMemoryProvider(ctx context.Context, options *Options) (*MemoryProvider, error) {
	driverName = MemoryDataProviderName
	dbHandle, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	if err = dbHandle.PingContext(ctx); err != nil {
		return nil, err
	}

	return &MemoryProvider{dbHandle: dbHandle}, nil
}

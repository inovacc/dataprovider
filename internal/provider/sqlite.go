package provider

import (
	"context"
	"fmt"
	"github.com/dyammarcano/dataprovider/internal/migration"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// SQLiteProvider defines the auth provider for SQLite database
type SQLiteProvider struct {
	dbHandle *sqlx.DB
}

// GetProviderStatus returns the status of the provider
func (s *SQLiteProvider) GetProviderStatus() Status {
	status := Status{
		Driver:   driverName,
		IsActive: true,
	}

	if err := s.CheckAvailability(); err != nil {
		status.IsActive = false
		status.Error = err
	}

	return status
}

// MigrateDatabase migrates the database to the latest version
func (s *SQLiteProvider) MigrateDatabase() migration.MigrationProvider {
	//TODO implement me
	panic("implement me")
}

// Disconnect disconnects from the data provider
func (s *SQLiteProvider) Disconnect() error {
	return s.dbHandle.Close()
}

// GetConnection returns the connection to the data provider
func (s *SQLiteProvider) GetConnection() *sqlx.DB {
	return s.dbHandle
}

// CheckAvailability checks if the data provider is available
func (s *SQLiteProvider) CheckAvailability() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()

	return s.dbHandle.PingContext(ctx)
}

// ReconnectDatabase reconnects to the database
func (s *SQLiteProvider) ReconnectDatabase() error {
	return s.CheckAvailability()
}

// InitializeDatabase initializes the database
func (s *SQLiteProvider) InitializeDatabase(schema string) error {
	_, err := s.dbHandle.Exec(schema)
	return err
}

// RevertDatabase reverts the database to the specified version
func (s *SQLiteProvider) RevertDatabase(targetVersion int) error {
	//TODO implement me
	panic("implement me")
}

// ResetDatabase resets the database
func (s *SQLiteProvider) ResetDatabase() error {
	//TODO implement me
	panic("implement me")
}

// NewSQLiteProvider creates a new SQLite provider instance
func NewSQLiteProvider(ctx context.Context, options *Options) (*SQLiteProvider, error) {
	driverName = SQLiteDataProviderName
	connectionString := options.ConnectionString

	if options.ConnectionString == "" {
		connectionString = fmt.Sprintf("file:%s.db?cache=shared&_foreign_keys=1", options.Name)
	}

	dbHandle, err := sqlx.Connect("sqlite", connectionString)
	if err != nil {
		return nil, err
	}

	dbHandle.SetMaxOpenConns(1)

	if err = dbHandle.PingContext(ctx); err != nil {
		return nil, err
	}

	return &SQLiteProvider{dbHandle: dbHandle}, nil
}

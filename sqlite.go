package dataprovider

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteProvider defines the auth provider for SQLite database
type SQLiteProvider struct {
	dbHandle *sqlx.DB
}

func (s *SQLiteProvider) GetProviderStatus() ProviderStatus {
	status := ProviderStatus{
		Driver:   driverName,
		IsActive: true,
	}

	if err := s.CheckAvailability(); err != nil {
		status.IsActive = false
		status.Error = err
	}

	return status
}

func (s *SQLiteProvider) MigrateDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (s *SQLiteProvider) Disconnect() error {
	return s.dbHandle.Close()
}

func (s *SQLiteProvider) GetConnection() *sqlx.DB {
	return s.dbHandle
}

func (s *SQLiteProvider) CheckAvailability() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()

	return s.dbHandle.PingContext(ctx)
}

func (s *SQLiteProvider) ReconnectDatabase() error {
	return s.CheckAvailability()
}

func (s *SQLiteProvider) InitializeDatabase(schema string) error {
	_, err := s.dbHandle.Exec(schema)
	return err
}

func (s *SQLiteProvider) RevertDatabase(targetVersion int) error {
	//TODO implement me
	panic("implement me")
}

func (s *SQLiteProvider) ResetDatabase() error {
	//TODO implement me
	panic("implement me")
}

func newSQLiteProvider(ctx context.Context, cfg *ConfigModule) (*SQLiteProvider, error) {
	connectionString := cfg.ConnectionString

	if cfg.ConnectionString == "" {
		connectionString = fmt.Sprintf("file:%s.db?cache=shared&_foreign_keys=1", cfg.Name)
	}

	dbHandle, err := sqlx.Connect("sqlite3", connectionString)
	if err != nil {
		return nil, err
	}

	dbHandle.SetMaxOpenConns(1)

	if err = dbHandle.PingContext(ctx); err != nil {
		return nil, err
	}

	return &SQLiteProvider{dbHandle: dbHandle}, nil
}

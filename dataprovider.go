package dataprovider

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	// OracleDatabaseProviderName defines the name for Oracle database Provider
	OracleDatabaseProviderName = "oracle"

	// SQLiteDataProviderName defines the name for SQLite database Provider
	SQLiteDataProviderName = "sqlite"

	// MySQLDatabaseProviderName defines the name for MySQL database Provider
	MySQLDatabaseProviderName = "mysql"

	// PostgreSQLDatabaseProviderName defines the name for PostgreSQL database Provider
	PostgreSQLDatabaseProviderName = "postgresql"

	// MemoryDataProviderName defines the name for memory provider using SQLite in-memory database Provider
	MemoryDataProviderName = "memory"
)

// ordering constants
const (
	OrderASC  = "ASC"
	OrderDESC = "DESC"
)

const (
	defaultSQLQueryTimeout = 10 * time.Second
	longSQLQueryTimeout    = 60 * time.Second
)

type Provider interface {
	// Disconnect disconnects from the data provider
	Disconnect() error

	// GetConnection returns the connection to the data provider
	GetConnection() *sqlx.DB

	// CheckAvailability checks if the data provider is available
	CheckAvailability() error

	// ReconnectDatabase reconnects to the database
	ReconnectDatabase() error

	// InitializeDatabase initializes the database
	InitializeDatabase(schema string) error

	// MigrateDatabase migrates the database to the latest version
	MigrateDatabase() error

	// RevertDatabase reverts the database to the specified version
	RevertDatabase(targetVersion int) error

	// ResetDatabase resets the database
	ResetDatabase() error
}

type ProviderStatus struct {
	Driver   string `json:"driver"`
	Error    error  `json:"error"`
	IsActive bool   `json:"is_active"`
}

type Wrapper struct {
	Driver  string
	Version int
	Provider
}

// NewProvider creates a new data provider instance
func NewProvider(ctx context.Context, cfg *ConfigModule) (*Wrapper, error) {
	switch cfg.Driver {
	case OracleDatabaseProviderName:
		return newOracleProvider(ctx, cfg)
	case SQLiteDataProviderName:
		return newSQLiteProvider(ctx, cfg)
	case MySQLDatabaseProviderName:
		return newMySQLProvider(ctx, cfg)
	case PostgreSQLDatabaseProviderName:
		return newPostgreSQLProvider(ctx, cfg)
	case MemoryDataProviderName:
		return newMemoryProvider(ctx, cfg)
	}

	return nil, fmt.Errorf("unsupported driver %s", cfg.Driver)
}

func (w *Wrapper) GetProviderStatus() ProviderStatus {
	status := ProviderStatus{
		Driver:   w.Driver,
		IsActive: true,
	}

	if err := w.CheckAvailability(); err != nil {
		status.IsActive = false
		status.Error = err
	}

	return status
}

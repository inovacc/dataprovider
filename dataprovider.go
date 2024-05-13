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

var driverName string

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

	// GetProviderStatus returns the status of the provider
	GetProviderStatus() ProviderStatus
}

type ProviderStatus struct {
	Driver   string `json:"driver"`
	Error    error  `json:"error"`
	IsActive bool   `json:"is_active"`
}

// NewDataProvider creates a new data provider instance
func NewDataProvider(ctx context.Context, cfg *ConfigModule) (Provider, error) {
	driverName = cfg.Driver

	switch driverName {
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

	return nil, fmt.Errorf("unsupported driver %s", driverName)
}

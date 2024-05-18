package dataprovider

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/afero"
	"path/filepath"
)

const (
	// OracleDatabaseProviderName defines the name for Oracle database Provider
	OracleDatabaseProviderName = "godror"

	// SQLiteDataProviderName defines the name for SQLite database Provider
	SQLiteDataProviderName = "sqlite"

	// MySQLDatabaseProviderName defines the name for MySQL database Provider
	MySQLDatabaseProviderName = "mysql"

	// PostgreSQLDatabaseProviderName defines the name for PostgreSQL database Provider
	PostgreSQLDatabaseProviderName = "postgres"

	// MemoryDataProviderName defines the name for memory provider using SQLite in-memory database Provider
	MemoryDataProviderName = "memory"
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
func NewDataProvider(cfg *ConfigModule) (Provider, error) {
	return NewDataProviderContext(context.Background(), cfg)
}

// NewDataProviderContext creates a new data provider instance
func NewDataProviderContext(ctx context.Context, cfg *ConfigModule) (Provider, error) {
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

func GetQueryFromFile(filename string) (string, error) {
	fs := afero.NewOsFs()

	ok, err := afero.DirExists(fs, filepath.Dir(filename))
	if err != nil {
		return "", err
	}

	if !ok {
		return "", fmt.Errorf("directory %s does not exist", filepath.Dir(filename))
	}

	ok, err = afero.Exists(fs, filename)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", fmt.Errorf("file %s does not exist", filename)
	}

	content, err := afero.ReadFile(fs, filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

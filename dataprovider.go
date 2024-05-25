package dataprovider

import (
	"fmt"
	"github.com/dyammarcano/dataprovider/internal/migration"
	"github.com/dyammarcano/dataprovider/internal/provider"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/afero"
	"path/filepath"
)

const (
	// OracleDatabaseProviderName defines the name for Oracle database Provider
	OracleDatabaseProviderName = provider.OracleDatabaseProviderName

	// SQLiteDataProviderName defines the name for SQLite database Provider
	SQLiteDataProviderName = provider.SQLiteDataProviderName

	// MySQLDatabaseProviderName defines the name for MySQL database Provider
	MySQLDatabaseProviderName = provider.MySQLDatabaseProviderName

	// PostgreSQLDatabaseProviderName defines the name for PostgreSQL database Provider
	PostgreSQLDatabaseProviderName = provider.PostgreSQLDatabaseProviderName

	// MemoryDataProviderName defines the name for memory provider using SQLite in-memory database Provider
	MemoryDataProviderName = provider.MemoryDataProviderName
)

type Status = provider.Status

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
	MigrateDatabase() migration.MigrationProvider

	// RevertDatabase reverts the database to the specified version
	RevertDatabase(targetVersion int) error

	// ResetDatabase resets the database
	ResetDatabase() error

	// GetProviderStatus returns the status of the provider
	GetProviderStatus() Status
}

// NewDataProvider creates a new data provider instance
func NewDataProvider(options *provider.Options) (Provider, error) {
	switch options.Driver {
	case OracleDatabaseProviderName:
		return provider.NewOracleProvider(options)
	case SQLiteDataProviderName:
		return provider.NewSQLiteProvider(options)
	case MySQLDatabaseProviderName:
		return provider.NewMySQLProvider(options)
	case PostgreSQLDatabaseProviderName:
		return provider.NewPostgreSQLProvider(options)
	case MemoryDataProviderName:
		return provider.NewMemoryProvider(options)
	}

	return nil, fmt.Errorf("unsupported driver %s", options.Driver)
}

// Must panics if the error is not nil
//
// Otherwise, it returns the provider instance with the corresponding implementation
func Must(provider Provider, err error) Provider {
	if err != nil {
		panic(err)
	}
	return provider
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

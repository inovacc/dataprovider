package dataprovider

import (
	"fmt"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/afero"
)

type databaseKind string

const (
	// OracleDatabaseProviderName defines the name for Oracle database Provider
	OracleDatabaseProviderName databaseKind = "oracle"

	// SQLiteDataProviderName defines the name for SQLite database Provider
	SQLiteDataProviderName databaseKind = "sqlite"

	// MySQLDatabaseProviderName defines the name for MySQL database Provider
	MySQLDatabaseProviderName databaseKind = "mysql"

	// PostgresSQLDatabaseProviderName defines the name for PostgresSQL database Provider
	PostgresSQLDatabaseProviderName databaseKind = "postgres"

	// MemoryDataProviderName defines the name for a memory provider using SQLite in-memory database Provider
	MemoryDataProviderName databaseKind = "memory"

	// SQLServerDatabaseProviderName defines the name for SQL Server database Provider
	SQLServerDatabaseProviderName databaseKind = "sqlserver"

	// MariadbDatabaseProviderName defines the name for MariaDB database Provider
	MariadbDatabaseProviderName databaseKind = "mariadb"
)

var driverName databaseKind

type Status struct {
	Driver   databaseKind `json:"driver"`
	Error    error        `json:"error"`
	IsActive bool         `json:"is_active"`
}

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
	MigrateDatabase() Migration

	// RevertDatabase reverts the database to the specified version
	RevertDatabase(targetVersion int) error

	// ResetDatabase resets the database
	ResetDatabase() error

	// GetProviderStatus returns the status of the provider
	GetProviderStatus() Status

	QueryBuilder() SQLBuilder
}

// NewDataProvider creates a new data provider instance
func NewDataProvider(options *Options) (Provider, error) {
	switch options.Driver {
	case OracleDatabaseProviderName:
		return NewOracleProvider(options)
	case SQLiteDataProviderName:
		return NewSQLiteProvider(options)
	case MySQLDatabaseProviderName:
		return NewMySQLProvider(options)
	case PostgresSQLDatabaseProviderName:
		return NewPostgreSQLProvider(options)
	case MemoryDataProviderName:
		return NewMemoryProvider(options)
	}

	return nil, fmt.Errorf("unsupported driver %s", options.Driver)
}

// Must launch panic if the error is not nil
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

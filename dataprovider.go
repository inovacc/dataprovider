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

// ConfigModule defines the configuration for the data provider
type ConfigModule struct {
	// Driver name, must be one of the SupportedProviders
	Driver string `json:"driver" mapstructure:"driver"`

	// Database name. For driver sqlite this can be the database name relative to the config dir
	// or the absolute path to the SQLite database.
	Name string `json:"name" mapstructure:"name"`

	// Database host. For postgresql and cockroachdb driver you can specify multiple hosts separated by commas
	Host string `json:"host" mapstructure:"host"`

	// Database port
	Port int `json:"port" mapstructure:"port"`

	// Database username
	Username     string `json:"username" mapstructure:"username"`
	UsernameFile string `json:"username_file" mapstructure:"username_file"`

	// Database password
	Password     string `json:"password" mapstructure:"password"`
	PasswordFile string `json:"password_file" mapstructure:"password_file"`

	// Database schema
	Schema string `json:"schema" mapstructure:"schema"`

	// prefix for SQL tables
	SQLTablesPrefix string `json:"sql_tables_prefix" mapstructure:"sql_tables_prefix"`

	// Sets the maximum number of open connections for mysql and postgresql driver.
	// Default 0 (unlimited)
	PoolSize int `json:"pool_size" mapstructure:"pool_size"`

	// Path to the backup directory. This can be an absolute path or a path relative to the config dir
	BackupsPath string `json:"backups_path" mapstructure:"backups_path"`

	// If not empty this connection string will be used instead of the other fields
	ConnectionString string `json:"connection_string" mapstructure:"connection_string"`
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

package provider

type NamedProvider string

const (
	// OracleDatabaseProviderName defines the name for Oracle database Provider
	OracleDatabaseProviderName NamedProvider = "oracle"

	// SQLiteDataProviderName defines the name for SQLite database Provider
	SQLiteDataProviderName NamedProvider = "sqlite"

	// MySQLDatabaseProviderName defines the name for MySQL database Provider
	MySQLDatabaseProviderName NamedProvider = "mysql"

	// PostgreSQLDatabaseProviderName defines the name for PostgreSQL database Provider
	PostgreSQLDatabaseProviderName NamedProvider = "postgres"

	// MemoryDataProviderName defines the name for memory provider using SQLite in-memory database Provider
	MemoryDataProviderName NamedProvider = "memory"
)

var driverName NamedProvider

type Status struct {
	Driver   NamedProvider `json:"driver"`
	Error    error         `json:"error"`
	IsActive bool          `json:"is_active"`
}

package dataprovider

const (
	// OracleDatabaseProviderName defines the name for Oracle database Provider
	OracleDatabaseProviderName string = "oracle"

	// SQLiteDataProviderName defines the name for SQLite database Provider
	SQLiteDataProviderName string = "sqlite"

	// MySQLDatabaseProviderName defines the name for MySQL database Provider
	MySQLDatabaseProviderName string = "mysql"

	// PostgresSQLDatabaseProviderName defines the name for PostgresSQL database Provider
	PostgresSQLDatabaseProviderName string = "postgres"

	// MemoryDataProviderName defines the name for a memory provider using SQLite in-memory database Provider
	MemoryDataProviderName string = "memory"

	// SQLServerDatabaseProviderName defines the name for SQL Server database Provider
	SQLServerDatabaseProviderName = "sqlserver"

	// MariadbDatabaseProviderName defines the name for MariaDB database Provider
	MariadbDatabaseProviderName = "mariadb"
)

var driverName string

type Status struct {
	Driver   string `json:"driver"`
	Error    error  `json:"error"`
	IsActive bool   `json:"is_active"`
}

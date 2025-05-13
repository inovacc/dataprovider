package dataprovider

import (
	"fmt"
	"strings"
)

type Dialect interface {
	QuoteIdentifier(name string) string
	AutoIncrement() string
	PrimaryKeySyntax() string
	TypeMap(goType string) string
	SupportsJSON() bool
	SupportsEnum() bool
	SupportsCheckConstraints() bool
	SupportsPartialIndex() bool
}

type PostgresDialect struct{}

func (d PostgresDialect) QuoteIdentifier(name string) string {
	return fmt.Sprintf(`"%s"`, name)
}

func (d PostgresDialect) AutoIncrement() string {
	return "SERIAL"
}

func (d PostgresDialect) PrimaryKeySyntax() string {
	return "PRIMARY KEY"
}

func (d PostgresDialect) TypeMap(goType string) string {
	switch goType {
	case "string":
		return "VARCHAR"
	case "int":
		return "INTEGER"
	case "bool":
		return "BOOLEAN"
	case "float64":
		return "DOUBLE PRECISION"
	case "time.Time":
		return "TIMESTAMP"
	default:
		return goType
	}
}

func (d PostgresDialect) SupportsJSON() bool {
	return true
}

func (d PostgresDialect) SupportsEnum() bool {
	return true
}

func (d PostgresDialect) SupportsCheckConstraints() bool {
	return true
}

func (d PostgresDialect) SupportsPartialIndex() bool {
	return true
}

type MySQLDialect struct{}

func (d MySQLDialect) QuoteIdentifier(name string) string {
	return fmt.Sprintf("`%s`", name)
}

func (d MySQLDialect) AutoIncrement() string {
	return "AUTO_INCREMENT"
}

func (d MySQLDialect) PrimaryKeySyntax() string {
	return "PRIMARY KEY"
}

func (d MySQLDialect) TypeMap(goType string) string {
	switch goType {
	case "string":
		return "VARCHAR(255)"
	case "int":
		return "INT"
	case "bool":
		return "TINYINT(1)"
	case "float64":
		return "DOUBLE"
	case "time.Time":
		return "DATETIME"
	default:
		return goType
	}
}

func (d MySQLDialect) SupportsJSON() bool {
	return true
}

func (d MySQLDialect) SupportsEnum() bool {
	return true
}

func (d MySQLDialect) SupportsCheckConstraints() bool {
	return false
}

func (d MySQLDialect) SupportsPartialIndex() bool {
	return false
}

type SQLiteDialect struct{}

func (d SQLiteDialect) SupportsJSON() bool {
	return false
}

func (d SQLiteDialect) SupportsEnum() bool {
	return false
}

func (d SQLiteDialect) SupportsCheckConstraints() bool {
	return true
}

func (d SQLiteDialect) SupportsPartialIndex() bool {
	return false
}

func (d SQLiteDialect) QuoteIdentifier(name string) string {
	return fmt.Sprintf(`"%s"`, name)
}

func (d SQLiteDialect) AutoIncrement() string {
	return "AUTOINCREMENT"
}

func (d SQLiteDialect) PrimaryKeySyntax() string {
	return "PRIMARY KEY"
}

func (d SQLiteDialect) TypeMap(goType string) string {
	switch goType {
	case "string":
		return "TEXT"
	case "int":
		return "INTEGER"
	case "bool":
		return "BOOLEAN"
	case "float64":
		return "REAL"
	case "time.Time":
		return "TEXT"
	default:
		return goType
	}
}

type OracleDialect struct{}

func (d OracleDialect) SupportsJSON() bool {
	return false
}

func (d OracleDialect) SupportsEnum() bool {
	return false
}

func (d OracleDialect) SupportsCheckConstraints() bool {
	return true
}

func (d OracleDialect) SupportsPartialIndex() bool {
	return false
}

func (d OracleDialect) QuoteIdentifier(name string) string {
	return fmt.Sprintf(`"%s"`, strings.ToUpper(name))
}

func (d OracleDialect) AutoIncrement() string {
	return "" // Se asume uso de SEQUENCE + TRIGGER externo
}

func (d OracleDialect) PrimaryKeySyntax() string {
	return "PRIMARY KEY"
}

func (d OracleDialect) TypeMap(goType string) string {
	switch goType {
	case "string":
		return "VARCHAR2(255)"
	case "int":
		return "NUMBER"
	case "bool":
		return "NUMBER(1)"
	case "float64":
		return "FLOAT"
	case "time.Time":
		return "TIMESTAMP"
	default:
		return goType
	}
}

type MariaDBDialect struct{}

func (d MariaDBDialect) SupportsJSON() bool {
	return true
}

func (d MariaDBDialect) SupportsEnum() bool {
	return true
}

func (d MariaDBDialect) SupportsCheckConstraints() bool {
	return false
}

func (d MariaDBDialect) SupportsPartialIndex() bool {
	return false
}

func (d MariaDBDialect) QuoteIdentifier(name string) string {
	return fmt.Sprintf("`%s`", name)
}

func (d MariaDBDialect) AutoIncrement() string {
	return "AUTO_INCREMENT"
}

func (d MariaDBDialect) PrimaryKeySyntax() string {
	return "PRIMARY KEY"
}

func (d MariaDBDialect) TypeMap(goType string) string {
	switch goType {
	case "string":
		return "VARCHAR(255)"
	case "int":
		return "INT"
	case "bool":
		return "TINYINT(1)"
	case "float64":
		return "DOUBLE"
	case "time.Time":
		return "DATETIME"
	default:
		return goType
	}
}

type SqlServerDialect struct{}

func (d SqlServerDialect) SupportsJSON() bool {
	return true
}

func (d SqlServerDialect) SupportsEnum() bool {
	return false
}

func (d SqlServerDialect) SupportsCheckConstraints() bool {
	return true
}

func (d SqlServerDialect) SupportsPartialIndex() bool {
	return false
}

func (d SqlServerDialect) QuoteIdentifier(name string) string {
	return fmt.Sprintf("[%s]", name)
}

func (d SqlServerDialect) AutoIncrement() string {
	return "IDENTITY(1,1)"
}

func (d SqlServerDialect) PrimaryKeySyntax() string {
	return "PRIMARY KEY"
}

func (d SqlServerDialect) TypeMap(goType string) string {
	switch goType {
	case "string":
		return "NVARCHAR(255)"
	case "int":
		return "INT"
	case "bool":
		return "BIT"
	case "float64":
		return "FLOAT"
	case "time.Time":
		return "DATETIME2"
	default:
		return goType
	}
}

func NewDialect(opts *Options) Dialect {
	switch opts.Driver {
	case MySQLDatabaseProviderName:
		return MySQLDialect{}
	case PostgresSQLDatabaseProviderName:
		return PostgresDialect{}
	case SQLiteDataProviderName, MemoryDataProviderName:
		return SQLiteDialect{}
	case OracleDatabaseProviderName:
		return OracleDialect{}
	case SQLServerDatabaseProviderName:
		return SqlServerDialect{}
	case MariadbDatabaseProviderName:
		return MariaDBDialect{}
	default:
		panic(fmt.Sprintf("unsupported dialect: %s", opts.Driver))
	}
}

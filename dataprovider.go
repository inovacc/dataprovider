package dataprovider

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/inovacc/dataprovider/internal/migration"
	"github.com/inovacc/dataprovider/internal/provider"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/spf13/afero"
	"path/filepath"
	"reflect"
)

const (
	// OracleDatabaseProviderName defines the name for Oracle database Provider
	OracleDatabaseProviderName = provider.OracleDatabaseProviderName

	// SQLiteDataProviderName defines the name for SQLite database Provider
	SQLiteDataProviderName = provider.SQLiteDataProviderName

	// MySQLDatabaseProviderName defines the name for MySQL database Provider
	MySQLDatabaseProviderName = provider.MySQLDatabaseProviderName

	// PostgresSQLDatabaseProviderName defines the name for PostgresSQL database Provider
	PostgresSQLDatabaseProviderName = provider.PostgresSQLDatabaseProviderName

	// MemoryDataProviderName defines the name for memory provider using SQLite in-memory database Provider
	MemoryDataProviderName = provider.MemoryDataProviderName
)

type Status = provider.Status
type Options = provider.Options

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
	MigrateDatabase() migration.Migration

	// RevertDatabase reverts the database to the specified version
	RevertDatabase(targetVersion int) error

	// ResetDatabase resets the database
	ResetDatabase() error

	// GetProviderStatus returns the status of the provider
	GetProviderStatus() Status

	// SqlBuilder returns the SQLBuilder instance
	SqlBuilder() *provider.SQLBuilder
}

// NewDataProvider creates a new data provider instance
func NewDataProvider(options *Options) (Provider, error) {
	switch options.Driver {
	case OracleDatabaseProviderName:
		return provider.NewOracleProvider(options)
	case SQLiteDataProviderName:
		return provider.NewSQLiteProvider(options)
	case MySQLDatabaseProviderName:
		return provider.NewMySQLProvider(options)
	case PostgresSQLDatabaseProviderName:
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

// StructScan scans a row from the database into a struct. The struct should be passed by reference. (riped out from sqlx)
func StructScan(rows *sql.Rows, dest any) error {
	v := reflect.ValueOf(dest)

	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}

	v = v.Elem()

	var (
		started bool
		fields  [][]int
		values  []any
		unsafe  bool
	)

	if !started {
		columns, err := rows.Columns()
		if err != nil {
			return err
		}
		m := &reflectx.Mapper{}

		fields = m.TraversalsByName(v.Type(), columns)
		// if we are not unsafe and are missing fields, return an error
		if f, err := missingFields(fields); err != nil && !unsafe {
			return fmt.Errorf("missing destination name %s in %T", columns[f], dest)
		}
		values = make([]any, len(columns))
		started = true
	}

	if err := fieldsByTraversal(v, fields, values, true); err != nil {
		return err
	}
	// scan into the struct field pointers and append to our results
	if err := rows.Scan(values...); err != nil {
		return err
	}
	return rows.Err()
}

func missingFields(transverse [][]int) (field int, err error) {
	for i, t := range transverse {
		if len(t) == 0 {
			return i, errors.New("missing field")
		}
	}
	return 0, nil
}

func fieldsByTraversal(v reflect.Value, traversals [][]int, values []any, ptrs bool) error {
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return errors.New("argument not a struct")
	}

	for i, traversal := range traversals {
		if len(traversal) == 0 {
			values[i] = new(any)
			continue
		}
		f := reflectx.FieldByIndexes(v, traversal)
		if ptrs {
			values[i] = f.Addr().Interface()
		} else {
			values[i] = f.Interface()
		}
	}
	return nil
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

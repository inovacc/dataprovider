package dataprovider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	Driver           databaseKind
	Name             string
	Host             string
	Port             int
	Username         string
	Password         string
	Schema           string
	SQLTablesPrefix  string
	PoolSize         int
	ConnectionString string
	context.Context
}

var opts *Options

func init() {
	opts = &Options{
		Context: context.Background(),
		Driver:  MemoryDataProviderName,
	}
}

type OptionFunc func(*Options)

// WithSqliteDB sets sqlite db path name
func WithSqliteDB(name, path string) OptionFunc {
	sb := strings.Builder{}

	if path == "." {
		dir, _ := os.Getwd()
		path = dir
	}

	sb.WriteString("file:")
	if !strings.HasSuffix(name, ".sqlite3") {
		sb.WriteString(fmt.Sprintf("%s.sqlite3", filepath.Join(path, name)))
	}

	sb.WriteString("?cache=shared")
	sb.WriteString("&mode=rwc")

	return func(o *Options) {
		o.ConnectionString = sb.String()
		o.Driver = SQLiteDataProviderName
	}
}

// WithMemoryDB sets memory db
func WithMemoryDB() OptionFunc {
	return func(o *Options) {
		o.ConnectionString = "file::memory:?cache=shared"
		o.Driver = MemoryDataProviderName
	}
}

// WithName sets db name
func WithName(name string) OptionFunc {
	return func(o *Options) {
		o.Name = name
	}
}

// WithDriver sets db driver
func WithDriver(driver databaseKind) OptionFunc {
	return func(o *Options) {
		o.Driver = driver
	}
}

// WithHost sets db host
func WithHost(host string) OptionFunc {
	return func(o *Options) {
		o.Host = host
	}
}

// WithPort sets db port
func WithPort(port int) OptionFunc {
	return func(o *Options) {
		o.Port = port
	}
}

// WithUsername sets db username
func WithUsername(username string) OptionFunc {
	return func(o *Options) {
		o.Username = username
	}
}

// WithPassword sets db password
func WithPassword(password string) OptionFunc {
	return func(o *Options) {
		o.Password = password
	}
}

// WithSQLTablesPrefix sets db sql tables prefix
func WithSQLTablesPrefix(sqlTablesPrefix string) OptionFunc {
	return func(o *Options) {
		o.SQLTablesPrefix = sqlTablesPrefix
	}
}

func WithPoolSize(poolSize int) OptionFunc {
	return func(o *Options) {
		o.PoolSize = poolSize
	}
}

// WithConnectionString sets db connection string
func WithConnectionString(connectionString string) OptionFunc {
	return func(o *Options) {
		o.ConnectionString = connectionString
	}
}

// WithContext sets db context
func WithContext(ctx context.Context) OptionFunc {
	return func(o *Options) {
		o.Context = ctx
	}
}

// NewOptions creates a new options instance
func NewOptions(optsFn ...OptionFunc) *Options {
	for _, opt := range optsFn {
		opt(opts)
	}
	return opts
}

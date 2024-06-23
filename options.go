package dataprovider

import (
	"context"
	"github.com/inovacc/dataprovider/internal/provider"
)

type OptionFunc func(*provider.Options)

func WithDriver(driver string) OptionFunc {
	return func(o *provider.Options) {
		o.Driver = driver
	}
}

func WithName(name string) OptionFunc {
	return func(o *provider.Options) {
		o.Name = name
	}
}

func WithHost(host string) OptionFunc {
	return func(o *provider.Options) {
		o.Host = host
	}
}

func WithPort(port int) OptionFunc {
	return func(o *provider.Options) {
		o.Port = port
	}
}

func WithUsername(username string) OptionFunc {
	return func(o *provider.Options) {
		o.Username = username
	}
}

func WithPassword(password string) OptionFunc {
	return func(o *provider.Options) {
		o.Password = password
	}
}

func WithSchema(schema string) OptionFunc {
	return func(o *provider.Options) {
		o.Schema = schema
	}
}

func WithSQLTablesPrefix(sqlTablesPrefix string) OptionFunc {
	return func(o *provider.Options) {
		o.SQLTablesPrefix = sqlTablesPrefix
	}
}

func WithPoolSize(poolSize int) OptionFunc {
	return func(o *provider.Options) {
		o.PoolSize = poolSize
	}
}

func WithConnectionString(connectionString string) OptionFunc {
	return func(o *provider.Options) {
		o.ConnectionString = connectionString
	}
}

func WithContext(ctx context.Context) OptionFunc {
	return func(o *provider.Options) {
		o.Context = ctx
	}
}

func defaultOptions() *provider.Options {
	return &provider.Options{
		Context: context.Background(),
		Driver:  MemoryDataProviderName,
	}
}

func NewOptions(opts ...OptionFunc) *provider.Options {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return o
}

//go:build !mysql

package dataprovider

import (
	"github.com/jmoiron/sqlx"
)

// MySQLProvider defines the auth provider for MySQL/MariaDB database
type MySQLProvider struct {
	options *Options
	dialect Dialect
}

func (m *MySQLProvider) QueryBuilder() SQLBuilder {
	return NewQueryBuilder(m.options)
}

func (m *MySQLProvider) Disconnect() error {
	// TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) GetConnection() *sqlx.DB {
	// TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) CheckAvailability() error {
	// TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) ReconnectDatabase() error {
	// TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) InitializeDatabase(_ string) error {
	// TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) MigrateDatabase() Migration {
	// TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) RevertDatabase(_ int) error {
	// TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) ResetDatabase() error {
	// TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) GetProviderStatus() Status {
	// TODO implement me
	panic("implement me")
}

// NewMySQLProvider creates a new MySQL provider instance
func NewMySQLProvider(_ *Options) (Provider, error) {
	panic("to use this driver you need to build with [mysql] tag")
}

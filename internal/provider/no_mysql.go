//go:build !mysql

package provider

import (
	"context"
	"github.com/dyammarcano/dataprovider/internal/migration"
	"github.com/jmoiron/sqlx"
)

// MySQLProvider defines the auth provider for MySQL/MariaDB database
type MySQLProvider struct {
	dbHandle *sqlx.DB
	context.Context
}

func (m *MySQLProvider) Disconnect() error {
	//TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) GetConnection() *sqlx.DB {
	//TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) CheckAvailability() error {
	//TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) ReconnectDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) InitializeDatabase(schema string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) MigrateDatabase() migration.MigrationProvider {
	//TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) RevertDatabase(targetVersion int) error {
	//TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) ResetDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) GetProviderStatus() Status {
	//TODO implement me
	panic("implement me")
}

// NewMySQLProvider creates a new MySQL provider instance
func NewMySQLProvider(options *Options) (*MySQLProvider, error) {
	panic("to use this driver you need to build with [mysql] tag")
}

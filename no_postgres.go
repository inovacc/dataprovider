//go:build !postgres

package dataprovider

import (
	"github.com/jmoiron/sqlx"
)

// PGSQLProvider defines the auth provider for PostgresSQL database
type PGSQLProvider struct {
	options *Options
	dialect Dialect
}

func (p *PGSQLProvider) QueryBuilder() SQLBuilder {
	return NewQueryBuilder(p.options)
}

func (p *PGSQLProvider) Disconnect() error {
	// TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) GetConnection() *sqlx.DB {
	// TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) CheckAvailability() error {
	// TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) ReconnectDatabase() error {
	// TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) InitializeDatabase(_ string) error {
	// TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) MigrateDatabase() Migration {
	// TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) RevertDatabase(_ int) error {
	// TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) ResetDatabase() error {
	// TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) GetProviderStatus() Status {
	// TODO implement me
	panic("implement me")
}

// NewPostgreSQLProvider creates a new PostgreSQL provider instance
func NewPostgreSQLProvider(_ *Options) (Provider, error) {
	panic("to use this driver you need to build with [postgres] tag")
}

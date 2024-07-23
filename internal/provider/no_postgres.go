//go:build !postgres

package provider

import (
	"github.com/inovacc/dataprovider/internal/migration"
	"github.com/jmoiron/sqlx"
)

// PGSQLProvider defines the auth provider for PostgresSQL database
type PGSQLProvider struct{}

func (p *PGSQLProvider) SqlBuilder() *SQLBuilder {
	//TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) Disconnect() error {
	//TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) GetConnection() *sqlx.DB {
	//TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) CheckAvailability() error {
	//TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) ReconnectDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) InitializeDatabase(schema string) error {
	//TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) MigrateDatabase() migration.Migration {
	//TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) RevertDatabase(targetVersion int) error {
	//TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) ResetDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) GetProviderStatus() Status {
	//TODO implement me
	panic("implement me")
}

// NewPostgreSQLProvider creates a new PostgreSQL provider instance
func NewPostgreSQLProvider(options *Options) (*PGSQLProvider, error) {
	panic("to use this driver you need to build with [postgres] tag")
}

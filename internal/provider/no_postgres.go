//go:build !postgres

package provider

import (
	"context"
	"github.com/dyammarcano/dataprovider/internal/migration"
	"github.com/jmoiron/sqlx"
)

// PGSQLProvider defines the auth provider for PostgreSQL database
type PGSQLProvider struct {
	dbHandle *sqlx.DB
	context.Context
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

func (p *PGSQLProvider) MigrateDatabase() migration.MigrationProvider {
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

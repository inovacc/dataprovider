//go:build !oracle

package provider

import (
	"github.com/inovacc/dataprovider/internal/migration"
	"github.com/jmoiron/sqlx"
)

// ORASQLProvider defines the auth provider for Oracle database
type ORASQLProvider struct{}

func (o *ORASQLProvider) SqlBuilder() *SQLBuilder {
	//TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) Disconnect() error {
	//TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) GetConnection() *sqlx.DB {
	//TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) CheckAvailability() error {
	//TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) ReconnectDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) InitializeDatabase(schema string) error {
	//TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) MigrateDatabase() migration.Migration {
	//TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) RevertDatabase(targetVersion int) error {
	//TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) ResetDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) GetProviderStatus() Status {
	//TODO implement me
	panic("implement me")
}

// NewOracleProvider creates a new Oracle provider instance
func NewOracleProvider(options *Options) (*ORASQLProvider, error) {
	panic("to use this driver you need to build with [oracle] tag")
}

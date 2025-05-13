//go:build !oracle

package dataprovider

import (
	"github.com/jmoiron/sqlx"
)

// ORASQLProvider defines the auth provider for an Oracle database
type ORASQLProvider struct {
	options *Options
	dialect Dialect
}

func (o *ORASQLProvider) QueryBuilder() SQLBuilder {
	return NewQueryBuilder(o.options)
}

func (o *ORASQLProvider) Disconnect() error {
	// TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) GetConnection() *sqlx.DB {
	// TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) CheckAvailability() error {
	// TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) ReconnectDatabase() error {
	// TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) InitializeDatabase(_ string) error {
	// TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) MigrateDatabase() Migration {
	// TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) RevertDatabase(_ int) error {
	// TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) ResetDatabase() error {
	// TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) GetProviderStatus() Status {
	// TODO implement me
	panic("implement me")
}

// NewOracleProvider creates a new Oracle provider instance
func NewOracleProvider(_ *Options) (Provider, error) {
	panic("to use this driver you need to build with [oracle] tag")
}

//go:build oracle

package provider

import (
	"context"
	"fmt"
	_ "github.com/godror/godror"
	"github.com/inovacc/dataprovider/internal/migration"
	"github.com/jmoiron/sqlx"
)

// ORASQLProvider defines the auth provider for Oracle database
type ORASQLProvider struct {
	dbHandle *sqlx.DB
	context.Context
}

func (o *ORASQLProvider) MigrateDatabase() migration.MigrationProvider {
	//TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) Disconnect() error {
	//TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) GetConnection() *sqlx.DB {
	return o.dbHandle
}

func (o *ORASQLProvider) CheckAvailability() error {
	ctx, cancel := context.WithTimeout(o.Context, 5)
	defer cancel()

	return o.dbHandle.PingContext(ctx)
}

func (o *ORASQLProvider) ReconnectDatabase() error {
	return o.CheckAvailability()
}

func (o *ORASQLProvider) InitializeDatabase(schema string) error {
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
	status := Status{
		Driver:   driverName,
		IsActive: true,
	}

	if err := o.CheckAvailability(); err != nil {
		status.IsActive = false
		status.Error = err
	}

	return status
}

// NewOracleProvider creates a new Oracle provider instance
func NewOracleProvider(options *Options) (*ORASQLProvider, error) {
	driverName = OracleDatabaseProviderName
	dataSourceName := fmt.Sprintf("%s/%s@%s:%d/%s",
		options.Username, options.Password, options.Host, options.Port, options.Name)

	dbHandle, err := sqlx.Connect("godror", dataSourceName)
	if err != nil {
		return nil, err
	}

	dbHandle.SetMaxOpenConns(options.PoolSize * 2)
	dbHandle.SetMaxIdleConns(options.PoolSize)

	if err = dbHandle.PingContext(options.Context); err != nil {
		return nil, err
	}

	return &ORASQLProvider{
		dbHandle: dbHandle,
		Context:  options.Context,
	}, nil
}

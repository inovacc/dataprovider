//go:build postgres

package dataprovider

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PGSQLProvider defines the auth provider for PostgresSQL database
type PGSQLProvider struct {
	dbHandle *sqlx.DB
	context.Context
}

func (p *PGSQLProvider) NewSQLBuilder() SQLBuilder {
	return NewQueryBuilder(p.options)
}

func (p *PGSQLProvider) GetProviderStatus() Status {
	status := Status{
		Driver:   driverName,
		IsActive: true,
	}

	if err := p.CheckAvailability(); err != nil {
		status.IsActive = false
		status.Error = err
	}

	return status
}

func (p *PGSQLProvider) MigrateDatabase() Migration {
	// TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) Disconnect() error {
	return p.dbHandle.Close()
}

func (p *PGSQLProvider) GetConnection() *sqlx.DB {
	return p.dbHandle
}

func (p *PGSQLProvider) CheckAvailability() error {
	ctx, cancel := context.WithTimeout(p.Context, 5)
	defer cancel()

	return p.dbHandle.PingContext(ctx)
}

func (p *PGSQLProvider) ReconnectDatabase() error {
	return p.CheckAvailability()
}

func (p *PGSQLProvider) InitializeDatabase(schema string) error {
	_, err := p.dbHandle.Exec(schema)
	return err
}

func (p *PGSQLProvider) RevertDatabase(targetVersion int) error {
	// TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) ResetDatabase() error {
	// TODO implement me
	panic("implement me")
}

// NewPostgresSQLProvider creates a new PostgresSQL provider instance
func NewPostgresSQLProvider(options *Options) (*PGSQLProvider, error) {
	driverName = options.Driver
	dataSourceName := fmt.Sprintf("user=%s dbname=%s password=%s port=%d host=%s sslmode=disable",
		options.Username, options.Name, options.Password, options.Port, options.Host)

	dbHandle, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	dbHandle.SetMaxOpenConns(options.PoolSize)
	if options.PoolSize > 0 {
		dbHandle.SetMaxIdleConns(options.PoolSize)
	} else {
		dbHandle.SetMaxIdleConns(2)
	}

	if err = dbHandle.PingContext(options.Context); err != nil {
		return nil, err
	}

	return &PGSQLProvider{
		dbHandle: dbHandle,
		Context:  options.Context,
		options:  options,
	}, nil
}

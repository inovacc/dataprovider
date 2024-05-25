package provider

import (
	"context"
	"fmt"
	"github.com/dyammarcano/dataprovider/internal/migration"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MySQLProvider defines the auth provider for MySQL/MariaDB database
type MySQLProvider struct {
	dbHandle *sqlx.DB
	context.Context
}

func (m *MySQLProvider) GetProviderStatus() Status {
	status := Status{
		Driver:   driverName,
		IsActive: true,
	}

	if err := m.CheckAvailability(); err != nil {
		status.IsActive = false
		status.Error = err
	}

	return status
}

func (m *MySQLProvider) MigrateDatabase() migration.MigrationProvider {
	//TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) Disconnect() error {
	return m.dbHandle.Close()
}

func (m *MySQLProvider) GetConnection() *sqlx.DB {
	return m.dbHandle
}

func (m *MySQLProvider) CheckAvailability() error {
	ctx, cancel := context.WithTimeout(m.Context, 5)
	defer cancel()

	return m.dbHandle.PingContext(ctx)
}

func (m *MySQLProvider) ReconnectDatabase() error {
	return m.CheckAvailability()
}

func (m *MySQLProvider) InitializeDatabase(schema string) error {
	_, err := m.dbHandle.Exec(schema)
	return err
}

func (m *MySQLProvider) RevertDatabase(targetVersion int) error {
	//TODO implement me
	panic("implement me")
}

func (m *MySQLProvider) ResetDatabase() error {
	//TODO implement me
	panic("implement me")
}

func NewMySQLProvider(options *Options) (*MySQLProvider, error) {
	driverName = MySQLDatabaseProviderName
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		options.Username, options.Password, options.Host, options.Port, options.Name)

	dbHandle, err := sqlx.Connect("mysql", dataSourceName)
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

	return &MySQLProvider{
		dbHandle: dbHandle,
		Context:  options.Context,
	}, nil
}

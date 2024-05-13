package dataprovider

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// MySQLProvider defines the auth provider for MySQL/MariaDB database
type MySQLProvider struct {
	dbHandle *sqlx.DB
}

func (m *MySQLProvider) MigrateDatabase() error {
	//TODO implement me
	panic("implement me")
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
	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()

	return m.dbHandle.PingContext(ctx)
}

func (m *MySQLProvider) ReconnectDatabase() error {
	return m.CheckAvailability()
}

func (m *MySQLProvider) InitializeDatabase(schema string) error {
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

func newMySQLProvider(ctx context.Context, cfg *ConfigModule) (*Wrapper, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	dbHandle, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	dbHandle.SetMaxOpenConns(cfg.PoolSize)
	if cfg.PoolSize > 0 {
		dbHandle.SetMaxIdleConns(cfg.PoolSize)
	} else {
		dbHandle.SetMaxIdleConns(2)
	}

	if err = dbHandle.PingContext(ctx); err != nil {
		return nil, err
	}

	return &Wrapper{
		Driver:   cfg.Driver,
		Version:  1,
		Provider: &MySQLProvider{dbHandle: dbHandle},
	}, nil
}

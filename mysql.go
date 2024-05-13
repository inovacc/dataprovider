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

func (m *MySQLProvider) InitializeDatabase() error {
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
	dsnString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	dbHandle, err := sqlx.Connect("mysql", dsnString)
	if err != nil {
		return nil, err
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

package dataprovider

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// ORASQLProvider defines the auth provider for Oracle database
type ORASQLProvider struct {
	dbHandle *sqlx.DB
}

func (o *ORASQLProvider) MigrateDatabase() error {
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
	ctx, cancel := context.WithTimeout(context.Background(), 5)
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

func newOracleProvider(ctx context.Context, cfg *ConfigModule) (*Wrapper, error) {
	dataSourceName := fmt.Sprintf("%s/%s@%s:%d/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	dbHandle, err := sqlx.Connect("godror", dataSourceName)
	if err != nil {
		return nil, err
	}

	dbHandle.SetMaxOpenConns(cfg.PoolSize * 2)
	dbHandle.SetMaxIdleConns(cfg.PoolSize)

	if err = dbHandle.PingContext(ctx); err != nil {
		return nil, err
	}

	return &Wrapper{
		Driver:   cfg.Driver,
		Version:  1,
		Provider: &ORASQLProvider{dbHandle: dbHandle},
	}, nil
}

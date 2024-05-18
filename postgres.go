package dataprovider

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PGSQLProvider defines the auth provider for PostgreSQL database
type PGSQLProvider struct {
	dbHandle *sqlx.DB
}

func (p *PGSQLProvider) GetProviderStatus() ProviderStatus {
	status := ProviderStatus{
		Driver:   driverName,
		IsActive: true,
	}

	if err := p.CheckAvailability(); err != nil {
		status.IsActive = false
		status.Error = err
	}

	return status
}

func (p *PGSQLProvider) MigrateDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) Disconnect() error {
	return p.dbHandle.Close()
}

func (p *PGSQLProvider) GetConnection() *sqlx.DB {
	return p.dbHandle
}

func (p *PGSQLProvider) CheckAvailability() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5)
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
	//TODO implement me
	panic("implement me")
}

func (p *PGSQLProvider) ResetDatabase() error {
	//TODO implement me
	panic("implement me")
}

func newPostgreSQLProvider(ctx context.Context, cfg *ConfigModule) (*PGSQLProvider, error) {
	dataSourceName := fmt.Sprintf("user=%s dbname=%s password=%s port=%d host=%s sslmode=disable", cfg.Username, cfg.Name, cfg.Password, cfg.Port, cfg.Host)
	dbHandle, err := sqlx.Connect(PostgreSQLDatabaseProviderName, dataSourceName)
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

	return &PGSQLProvider{dbHandle: dbHandle}, nil
}

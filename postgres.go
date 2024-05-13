package dataprovider

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// PGSQLProvider defines the auth provider for PostgreSQL database
type PGSQLProvider struct {
	dbHandle *sqlx.DB
}

func (p *PGSQLProvider) MigrateDatabase() error {
	//TODO implement me
	panic("implement me")
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
	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()

	return p.dbHandle.PingContext(ctx)
}

func (p *PGSQLProvider) ReconnectDatabase() error {
	return p.CheckAvailability()
}

func (p *PGSQLProvider) InitializeDatabase() error {
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

func newPostgreSQLProvider(ctx context.Context, cfg *ConfigModule) (*Wrapper, error) {
	dsnString := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", cfg.Username, cfg.Name, cfg.Password)
	dbHandle, err := sqlx.Connect("postgres", dsnString)
	if err != nil {
		return nil, err
	}

	if err = dbHandle.PingContext(ctx); err != nil {
		return nil, err
	}

	return &Wrapper{
		Driver:   cfg.Driver,
		Version:  1,
		Provider: &PGSQLProvider{dbHandle: dbHandle},
	}, nil
}

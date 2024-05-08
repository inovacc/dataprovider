package dataprovider

import (
	"context"
	"github.com/jmoiron/sqlx"
)

// PGSQLProvider defines the auth provider for PostgreSQL database
type PGSQLProvider struct {
	dbHandle *sqlx.DB
}

func (P *PGSQLProvider) Connect() error {
	//TODO implement me
	panic("implement me")
}

func (P *PGSQLProvider) Disconnect() error {
	//TODO implement me
	panic("implement me")
}

func (P *PGSQLProvider) GetConnection() *sqlx.DB {
	//TODO implement me
	panic("implement me")
}

func (P *PGSQLProvider) CheckAvailability() error {
	//TODO implement me
	panic("implement me")
}

func (P *PGSQLProvider) ReconnectDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (P *PGSQLProvider) InitializeDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (P *PGSQLProvider) migrateDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (P *PGSQLProvider) RevertDatabase(targetVersion int) error {
	//TODO implement me
	panic("implement me")
}

func (P *PGSQLProvider) resetDatabase() error {
	//TODO implement me
	panic("implement me")
}

func newPostgreSQLProvider(ctx context.Context) error {
	dbHandle, err := sqlx.Connect("postgres", "user=postgres dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		return err
	}

	if err = dbHandle.PingContext(ctx); err != nil {
		return err
	}

	provider = &PGSQLProvider{dbHandle: dbHandle}

	return nil
}

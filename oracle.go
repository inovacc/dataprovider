package dataprovider

import (
	"context"
	"github.com/jmoiron/sqlx"
)

// ORASQLProvider defines the auth provider for Oracle database
type ORASQLProvider struct {
	dbHandle *sqlx.DB
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
	//TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) ReconnectDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) InitializeDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (o *ORASQLProvider) migrateDatabase() error {
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

func newOracleProvider(ctx context.Context) error {
	dbHandle, err := sqlx.ConnectContext(ctx, OracleDatabaseProviderName, "")
	if err != nil {
		return err
	}

	provider = &ORASQLProvider{dbHandle: dbHandle}

	return nil
}

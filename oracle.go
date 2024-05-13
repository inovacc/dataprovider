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
	ctxValue := ctx.Value("config").(*ConfigModule)
	if ctxValue == nil {
		return fmt.Errorf("config not found in context")
	}

	dsnString := fmt.Sprintf("%s/%s@%s:%d/%s", ctxValue.Username, ctxValue.Password, ctxValue.Host, ctxValue.Port, ctxValue.Name)

	dbHandle, err := sqlx.Connect("godror", dsnString)
	if err != nil {
		return err
	}

	provider = &ORASQLProvider{dbHandle: dbHandle}

	return nil
}

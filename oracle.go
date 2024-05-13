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

func (o *ORASQLProvider) InitializeDatabase() error {
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

func newOracleProvider(ctx context.Context) (*Wrapper, error) {
	ctxValue := ctx.Value("config").(*ConfigModule)
	if ctxValue == nil {
		return nil, fmt.Errorf("config not found in context")
	}

	dsnString := fmt.Sprintf("%s/%s@%s:%d/%s", ctxValue.Username, ctxValue.Password, ctxValue.Host, ctxValue.Port, ctxValue.Name)

	dbHandle, err := sqlx.Connect("godror", dsnString)
	if err != nil {
		return nil, err
	}

	return &Wrapper{
		Version:  1,
		Provider: &ORASQLProvider{dbHandle: dbHandle},
	}, nil
}

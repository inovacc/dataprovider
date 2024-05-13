package dataprovider

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// SQLiteProvider defines the auth provider for SQLite database
type SQLiteProvider struct {
	dbHandle *sqlx.DB
}

func (S *SQLiteProvider) Disconnect() error {
	//TODO implement me
	panic("implement me")
}

func (S *SQLiteProvider) GetConnection() *sqlx.DB {
	//TODO implement me
	panic("implement me")
}

func (S *SQLiteProvider) CheckAvailability() error {
	//TODO implement me
	panic("implement me")
}

func (S *SQLiteProvider) ReconnectDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (S *SQLiteProvider) InitializeDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (S *SQLiteProvider) migrateDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (S *SQLiteProvider) RevertDatabase(targetVersion int) error {
	//TODO implement me
	panic("implement me")
}

func (S *SQLiteProvider) ResetDatabase() error {
	//TODO implement me
	panic("implement me")
}

func newSQLiteProvider(ctx context.Context) error {
	ctxValue := ctx.Value("config").(*ConfigModule)
	if ctxValue == nil {
		return fmt.Errorf("config not found in context")
	}

	dbHandle, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		return err
	}

	if err = dbHandle.PingContext(ctx); err != nil {
		return err
	}

	provider = &SQLiteProvider{dbHandle: dbHandle}

	return nil
}

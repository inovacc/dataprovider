package dataprovider

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteProvider defines the auth provider for SQLite database
type SQLiteProvider struct {
	dbHandle *sqlx.DB
}

func (s *SQLiteProvider) MigrateDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (s *SQLiteProvider) Disconnect() error {
	//TODO implement me
	panic("implement me")
}

func (s *SQLiteProvider) GetConnection() *sqlx.DB {
	//TODO implement me
	panic("implement me")
}

func (s *SQLiteProvider) CheckAvailability() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()

	return s.dbHandle.PingContext(ctx)
}

func (s *SQLiteProvider) ReconnectDatabase() error {
	return s.CheckAvailability()
}

func (s *SQLiteProvider) InitializeDatabase() error {
	//TODO implement me
	panic("implement me")
}

func (s *SQLiteProvider) RevertDatabase(targetVersion int) error {
	//TODO implement me
	panic("implement me")
}

func (s *SQLiteProvider) ResetDatabase() error {
	//TODO implement me
	panic("implement me")
}

func newSQLiteProvider(ctx context.Context) (*Wrapper, error) {
	ctxValue := ctx.Value("config").(*ConfigModule)
	if ctxValue == nil {
		return nil, fmt.Errorf("config not found in context")
	}

	dbHandle, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	if err = dbHandle.PingContext(ctx); err != nil {
		return nil, err
	}

	return &Wrapper{
		Version:  1,
		Provider: &SQLiteProvider{dbHandle: dbHandle},
	}, nil
}

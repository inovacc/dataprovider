package dataprovider

import "github.com/jmoiron/sqlx"

// SQLiteProvider defines the auth provider for SQLite database
type SQLiteProvider struct {
	dbHandle *sqlx.DB
}

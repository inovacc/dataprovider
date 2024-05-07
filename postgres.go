package dataprovider

import "github.com/jmoiron/sqlx"

// PGSQLProvider defines the auth provider for PostgreSQL database
type PGSQLProvider struct {
	dbHandle *sqlx.DB
}

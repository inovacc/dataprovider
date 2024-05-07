package dataprovider

import "github.com/jmoiron/sqlx"

// ORASQLProvider defines the auth provider for Oracle database
type ORASQLProvider struct {
	dbHandle *sqlx.DB
}

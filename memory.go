package dataprovider

import "github.com/jmoiron/sqlx"

// MemoryProvider defines the auth provider for in-memory database
type MemoryProvider struct {
	dbHandle *sqlx.DB
}

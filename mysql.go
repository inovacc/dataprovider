package dataprovider

import "github.com/jmoiron/sqlx"

// MySQLProvider defines the auth provider for MySQL/MariaDB database
type MySQLProvider struct {
	dbHandle *sqlx.DB
}

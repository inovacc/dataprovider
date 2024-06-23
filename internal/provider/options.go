package provider

import "context"

type Options struct {
	Driver           string
	Name             string
	Host             string
	Port             int
	Username         string
	Password         string
	Schema           string
	SQLTablesPrefix  string
	PoolSize         int
	ConnectionString string
	context.Context
}

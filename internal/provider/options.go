package provider

type Options struct {
	Driver           NamedProvider
	Name             string
	Host             string
	Port             int
	Username         string
	Password         string
	Schema           string
	SQLTablesPrefix  string
	PoolSize         int
	ConnectionString string
}

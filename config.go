package dataprovider

// ConfigModule defines the configuration for the data provider
type ConfigModule struct {
	// Driver name, must be one of the SupportedProviders
	Driver string `json:"driver" mapstructure:"driver"`

	// Database name. For driver sqlite this can be the database name relative to the config dir
	// or the absolute path to the SQLite database.
	Name string `json:"name" mapstructure:"name"`

	// Database host. For postgresql and cockroachdb driver you can specify multiple hosts separated by commas
	Host string `json:"host" mapstructure:"host"`

	// Database port
	Port int `json:"port" mapstructure:"port"`

	// Database username
	Username     string `json:"username" mapstructure:"username"`
	UsernameFile string `json:"username_file" mapstructure:"username_file"`

	// Database password
	Password     string `json:"password" mapstructure:"password"`
	PasswordFile string `json:"password_file" mapstructure:"password_file"`

	// Database schema
	Schema string `json:"schema" mapstructure:"schema"`

	// prefix for SQL tables
	SQLTablesPrefix string `json:"sql_tables_prefix" mapstructure:"sql_tables_prefix"`

	// Sets the maximum number of open connections for mysql and postgresql driver.
	// Default 0 (unlimited)
	PoolSize int `json:"pool_size" mapstructure:"pool_size"`

	// Path to the backup directory. This can be an absolute path or a path relative to the config dir
	BackupsPath string `json:"backups_path" mapstructure:"backups_path"`

	// If not empty this connection string will be used instead of the other fields
	ConnectionString string `json:"connection_string" mapstructure:"connection_string"`
}

type Builder struct {
	driver           string
	name             string
	host             string
	port             int
	username         string
	password         string
	schema           string
	sqlTablesPrefix  string
	poolSize         int
	connectionString string
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) WithDriver(driver string) *Builder {
	b.driver = driver
	return b
}

func (b *Builder) WithName(name string) *Builder {
	b.name = name
	return b
}

func (b *Builder) WithHost(host string) *Builder {
	b.host = host
	return b
}

func (b *Builder) WithPort(port int) *Builder {
	b.port = port
	return b
}

func (b *Builder) WithUsername(username string) *Builder {
	b.username = username
	return b
}

func (b *Builder) WithPassword(password string) *Builder {
	b.password = password
	return b
}

func (b *Builder) WithSchema(schema string) *Builder {
	b.schema = schema
	return b
}

func (b *Builder) WithSQLTablesPrefix(sqlTablesPrefix string) *Builder {
	b.sqlTablesPrefix = sqlTablesPrefix
	return b
}

func (b *Builder) WithPoolSize(poolSize int) *Builder {
	b.poolSize = poolSize
	return b
}

func (b *Builder) WithConnectionString(connectionString string) *Builder {
	b.connectionString = connectionString
	return b
}

func (b *Builder) Build() *ConfigModule {
	return &ConfigModule{
		Driver:           b.driver,
		Name:             b.name,
		Host:             b.host,
		Port:             b.port,
		Username:         b.username,
		Password:         b.password,
		Schema:           b.schema,
		SQLTablesPrefix:  b.sqlTablesPrefix,
		PoolSize:         b.poolSize,
		ConnectionString: b.connectionString,
	}
}

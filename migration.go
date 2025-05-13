package dataprovider

type Migration interface {
	Validate(string) error
	Migrate() error
	Revert() error
}

type migrationProvider struct{}

func (m *migrationProvider) Validate(path string) error {
	// TODO implement me
	panic("implement me")
}

func (m *migrationProvider) Migrate() error {
	// TODO implement me
	panic("implement me")
}

func (m *migrationProvider) Revert() error {
	// TODO implement me
	panic("implement me")
}

func NewMigration() Migration {
	return &migrationProvider{}
}

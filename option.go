package ezdb

type Option func(*DB)

func WithName(name string) Option {
	return func(db *DB) {
		db.Name = name
	}
}

func WithConnector(c Connector) Option {
	return func(db *DB) {
		db.connector = c
	}
}

func WithMigrationDir(dir string) Option {
	return func(db *DB) {
		db.migrationDir = dir
	}
}

func WithDownMigrations(enabled bool) Option {
	return func(db *DB) {
		db.downEnabled = enabled
	}
}

// WithSchema pins the schema (search_path) used for the connection, including
// the schema that golang-migrate reads/writes its schema_migrations table in.
// Use this when the login role's default search_path does not point at the
// schema that holds the application's objects.
func WithSchema(name string) Option {
	return func(db *DB) {
		db.schema = name
	}
}

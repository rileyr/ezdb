package ezdb

type Option func(*DB)

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

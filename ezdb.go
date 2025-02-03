package ezdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// DB implements the functionality provided by the ezdb package
type DB struct {
	connector    Connector
	migrationDir string
	downEnabled  bool
}

// New returns a new DB with the given options.
func New(opts ...Option) *DB {
	d := &DB{
		downEnabled: true,
	}

	for _, opt := range opts {
		opt(d)
	}

	if d.connector == nil {
		d.connector = NewDefaultEnvConnector()
	}
	if d.migrationDir == "" {
		d.migrationDir = "./db/migrations"
	}

	return d
}

// CreateDatabase creates the database, for initial environment setup.
func (db *DB) CreateDatabase() error {
	connData, err := db.connector.Data()
	if err != nil {
		return err
	}

	// hack to dodge PGDATABASE if set:
	existing := os.Getenv("PGDATABASE")
	os.Setenv("PGDATABASE", "")
	defer func() {
		os.Setenv("PGDATABASE", existing)
	}()

	conn, err := sql.Open("postgres", connData.DSN())
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Exec("CREATE DATABASE " + connData.Db + ";")
	return err
}

// CreateMigration creates set of migration files with the current timestamp and given name.
func (db *DB) CreateMigration(name string) error {
	if name == "" {
		return errors.New("cannot create migration without name")
	}

	files, err := os.ReadDir(db.migrationDir)
	if err != nil {
		return fmt.Errorf("failed to prepare filename: %w", err)
	}

	fileVersions := []int{}
	for _, file := range files {
		n := file.Name()
		if strings.Contains(n, ".sql") && strings.Contains(n, "_") { // i am not smart enuf for a regex
			parts := strings.Split(n, "_")
			if num, err := strconv.Atoi(parts[0]); err == nil {
				fileVersions = append(fileVersions, num)
			}
		}
	}

	sort.SliceStable(fileVersions, func(i, j int) bool {
		return fileVersions[i] < fileVersions[j]
	})

	var newVersion int
	if len(fileVersions) > 0 {
		newVersion = fileVersions[len(fileVersions)-1] + 1
	} else {
		newVersion = 1
	}
	migrationName := fmt.Sprintf("%04d_%s", newVersion, name)

	create := func(filepath string) error {
		fmt.Printf("creating file: %s\n", filepath)
		_, err := os.Create(filepath)
		return err
	}

	path := db.migrationDir + "/" + migrationName
	if err := create(fmt.Sprintf("%s.up.sql", path)); err != nil {
		return err
	}

	if db.downEnabled {
		if err := create(fmt.Sprintf("%s.down.sql", path)); err != nil {
			return err
		}
	}

	return nil
}

// MigrateAll applies all outstanding migrations to the database.
func (db *DB) MigrateAll() error {
	mig, err := db.buildMigrator()
	if err != nil {
		return err
	}

	return db.migrateWith(mig, mig.Up)
}

// MigrateSteps applies a specific number of migrations to the database. If
// steps are negative, migrations will be rolled back.
func (db *DB) MigrateSteps(steps int) error {
	mig, err := db.buildMigrator()
	if err != nil {
		return err
	}

	return db.migrateWith(mig, func() error { return mig.Steps(steps) })
}

func (db *DB) migrateWith(mig *migrate.Migrate, do func() error) error {
	currentDbVersion, _, err := mig.Version()
	if err != nil && err.Error() != "no migration" {
		return err
	}

	err = do()
	if err != nil && err.Error() != "no change" {
		if err.Error() == "file does not exist" {
			return errors.New("found no applicable migrations")
		} else {
			return errors.Wrap(err, "migrate error")
		}
	}

	newDbVersion, _, err := mig.Version()
	if err != nil && err.Error() != "no migration" {
		return err
	}

	log.Printf("migrated %d to %d\n", currentDbVersion, newDbVersion)
	return nil
}

func (db *DB) buildMigrator() (*migrate.Migrate, error) {
	connData, err := db.connector.Data()
	if err != nil {
		return nil, err
	}

	conn, err := sql.Open("postgres", connData.ConnString())
	if err != nil {
		return nil, errors.Wrap(err, "open db error")
	}

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "open driver error")
	}

	fp := fmt.Sprintf("file://%s", db.migrationDir)
	return migrate.NewWithDatabaseInstance(fp, "postgres", driver)
}

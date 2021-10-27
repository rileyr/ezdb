# ezdb

EZDB (easy database) provides lightweight, simple database tooling for go applications, inspired by Rails. The specific featureset is:

  - create databases
  - manage migration files
  - apply / rollback migrations

## Quick Start

```golang
package main

import (
  "github.com/rileyr/ezdb"
  "log"
)


func main() {
  db := ezdb.New()

  // Create the database for the first time:
  if err := db.CreateDatabase(); err != nil {
    log.Fatal(err)
  }

  // Create a new migration:
  if err := db.CreateMigration("create_some_new_table"); err != nil {
    log.Fatal(err)
  }

  // Manually edit the migration file...

  // Apply all pending migrations:
  if err := db.MigrateAll(); err != nil {
    log.Fatal(err)
  }

  // Apply exactly one migration:
  if err := db.MigrateSteps(1); err != nil {
    log.Fatal(err)
  }

  // Roll back one migration:
  if err := db.MigrateSteps(-1); err != nil {
    log.Fatal(err)
  }
}
```

## Connection Details

By default, EZDB uses the default postgres environment variables:

  - `PGUSER` - database username
  - `PGPASSWORD` - database password
  - `PGHOST` - database host
  - `PGDATABASE` - database name
  - `PGPORT` - database port

## CLI

You can use the entrypoint defined in `./cmd/cli/main.go` as a standalone entrypoint, or, you can include it into your application's CLI:

```golang
package main

import(
  "github.com/spf13/cobra"
  "github.com/rileyr/ezdb/cmd"
)

func myAppCLI() *cobra.Command {
  c := &cobra.Command{}
  // blah blah blah
  return c
}

func main() {
  c := myAppCLI()
  c.AddCommand(cmd.NewCommand())
  c.Execute()
}
```

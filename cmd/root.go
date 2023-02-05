package cmd

import (
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/rileyr/ezdb"
	"github.com/spf13/cobra"
)

func newRootCommand(opts ...ezdb.Option) *cobra.Command {
	c := &cobra.Command{
		Use: "db",
		Long: `
	ezdb database operations

	Set connection information in the environment via the PG* env vars.
		`,
		PersistentPreRunE: createDB(opts...),
	}

	c.PersistentFlags().StringVarP(&migrationDir, "migrations", "m", "./db/migrations", "path to migrations dir")
	return c
}

var (
	db           *ezdb.DB
	migrationDir string
)

func createDB(opts ...ezdb.Option) func(c *cobra.Command, args []string) error {
	return func(c *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "get wd error")
		}
		migrationDir := path.Join(wd, migrationDir)

		opts = append(opts, ezdb.WithMigrationDir(migrationDir))
		db = ezdb.New(opts...)
		return nil
	}
}

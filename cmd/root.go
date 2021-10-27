package cmd

import (
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/rileyr/ezdb"
	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	c := &cobra.Command{
		Use: "db",
		Long: `
	ezdb database operations

	Set connection information in the environment via the PG* env vars.
		`,
		PersistentPreRunE: createDB,
	}

	c.PersistentFlags().StringVarP(&migrationDir, "migrations", "m", "./db/migrations", "path to migrations dir")
	return c
}

var (
	db           *ezdb.DB
	migrationDir string
)

func createDB(c *cobra.Command, args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "get wd error")
	}
	migrationDir := path.Join(wd, migrationDir)

	db = ezdb.New(ezdb.WithMigrationDir(migrationDir))
	return nil
}

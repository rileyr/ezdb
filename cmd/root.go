package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/rileyr/ezdb"
	"github.com/spf13/cobra"
)

func newRootCommand(cmdName string, opts ...ezdb.Option) *cobra.Command {
	c := &cobra.Command{
		Use: cmdName,
		Long: `
	ezdb database operations

	Set connection information in the environment via the PG* env vars.
		`,
	}

	c.PersistentFlags().StringVarP(&migrationDir, "migrations", "m", "./db/migrations", "path to migrations dir")
	return c
}

var (
	db           *ezdb.DB
	migrationDir string
)

func setupEzdbInstance(
	preRunE func(*cobra.Command, []string) error,
	dynamicOpts func() ([]ezdb.Option, error),
	staticOpts ...ezdb.Option,
) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, args []string) error {
		// If the original base command has a PRE configured, call it first:
		if preRunE != nil {
			if err := preRunE(c, args); err != nil {
				return err
			}
		}

		wd, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "get wd error")
		}
		migrationDir := path.Join(wd, migrationDir)

		opts := append(staticOpts, ezdb.WithMigrationDir(migrationDir))
		if dynamicOpts != nil {
			moreOpts, err := dynamicOpts()
			if err != nil {
				return fmt.Errorf("dynamic opts error: %w", err)
			}
			opts = append(opts, moreOpts...)
		}
		db = ezdb.New(opts...)
		return nil
	}
}

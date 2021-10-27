package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

func newCreateMigrationCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "new [NAME]",
		Short: "create a new migration",
		RunE:  runNew,
	}
}

func runNew(c *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("new migration requires a name")
	}

	return db.CreateMigration(args[0])
}

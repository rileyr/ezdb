package cmd

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	c := newRootCommand()
	c.AddCommand(newCreateCommand())
	c.AddCommand(newCreateMigrationCommand())
	c.AddCommand(newMigrateCommand())
	return c
}

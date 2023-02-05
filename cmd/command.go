package cmd

import (
	"github.com/rileyr/ezdb"
	"github.com/spf13/cobra"
)

func NewCommand(opts ...ezdb.Option) *cobra.Command {
	c := newRootCommand(opts...)
	c.AddCommand(newCreateCommand())
	c.AddCommand(newCreateMigrationCommand())
	c.AddCommand(newMigrateCommand())
	return c
}

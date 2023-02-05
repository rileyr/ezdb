package cmd

import (
	"github.com/rileyr/ezdb"
	"github.com/spf13/cobra"
)

func NewCommand(opts ...ezdb.Option) *cobra.Command {
	c := newRootCommand(opts...)

	ConfigureCommand(c, opts...)

	return c
}

func ConfigureCommand(c *cobra.Command, opts ...ezdb.Option) {
	c.PersistentPreRunE = setupEzdbInstance(c.PersistentPreRunE, opts...)
	c.AddCommand(newCreateCommand())
	c.AddCommand(newCreateMigrationCommand())
	c.AddCommand(newMigrateCommand())
}

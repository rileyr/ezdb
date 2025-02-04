package cmd

import (
	"github.com/rileyr/ezdb"
	"github.com/spf13/cobra"
)

func NewCommand(cmdName string, opts ...ezdb.Option) *cobra.Command {
	c := newRootCommand(cmdName, opts...)

	ConfigureCommand(c, nil, opts...)

	return c
}

func ConfigureCommand(
	c *cobra.Command,
	dynamicOpts func() ([]ezdb.Option, error),
	staticOpts ...ezdb.Option,
) {
	c.PersistentPreRunE = setupEzdbInstance(c.PersistentPreRunE, dynamicOpts, staticOpts...)
	c.AddCommand(newCreateCommand())
	c.AddCommand(newCreateMigrationCommand())
	c.AddCommand(newMigrateCommand())
}

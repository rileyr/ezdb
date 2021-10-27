package cmd

import "github.com/spf13/cobra"

func newCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "create a new database",
		RunE:  runCreate,
	}
}

func runCreate(c *cobra.Command, args []string) error {
	return db.CreateDatabase()
}

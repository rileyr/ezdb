package cmd

import "github.com/spf13/cobra"

func newMigrateCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "migrate",
		Short: "apply migrations to the database",
		RunE:  runMigrate,
	}

	c.Flags().IntVarP(&migrateSteps, "steps", "s", 0, "steps to apply(0=all,negative=rollback)")

	return c
}

var (
	migrateSteps int
)

func runMigrate(c *cobra.Command, args []string) error {
	if migrateSteps == 0 {
		return db.MigrateAll()
	} else {
		return db.MigrateSteps(migrateSteps)
	}
}

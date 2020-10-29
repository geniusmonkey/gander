package cmd

import (
	"github.com/geniusmonkey/gander/db"
	"github.com/geniusmonkey/gander/migration"
	"github.com/spf13/cobra"
	"log"
)

var baselineVer int64

var baselineCmd = &cobra.Command{
	Use:     "baseline VERSION",
	Short:   "Baseline an existing db to a specific VERSION",
	PreRun:  setup,
	PostRun: tearDown,
	Run: func(cmd *cobra.Command, args []string) {
		if baselineVer == 0 {
			log.Fatalf("you must supply a --version flag")
		}

		if err := migration.Baseline(db.Get(), proj.MigrationDir(), baselineVer); err != nil {
			log.Fatalf("failed to create migration, %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(baselineCmd)
	baselineCmd.Flags().Int64Var(&baselineVer, "version", 0, "version to migrate to")
}

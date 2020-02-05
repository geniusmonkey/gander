package cmd

import (
	"github.com/geniusmonkey/gander"
	"github.com/geniusmonkey/gander/db"
	"github.com/spf13/cobra"
	"log"
)

var redoCmd = &cobra.Command{
	Use: "redo",
	Short: "Re-run the latest migration",
	PreRun: setup,
	PostRun: tearDown,
	Run: func(cmd *cobra.Command, args []string) {
		if err := gander.Redo(db.Get(), proj.MigrationDir()); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(redoCmd)
}
package main

import (
	"log"
	"os"

	"github.com/autopticio/instance/src/entity"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "autopticli"}

	// Register entity commands
	rootCmd.AddCommand(entity.InventoryCommand())
	rootCmd.AddCommand(entity.UICommand())
	rootCmd.AddCommand(entity.StorybookCommand())
	rootCmd.AddCommand(entity.APICommand())
	rootCmd.AddCommand(entity.SchedulerCommand())
	rootCmd.AddCommand(entity.ModelCommand())

	// Execute the CLI
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

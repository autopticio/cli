package entity

import "github.com/spf13/cobra"

// Model Entity Commands
func ModelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "model",
		Short: "Commands related to the model and model service",
	}
	return cmd
}

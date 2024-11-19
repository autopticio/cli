package entity

import "github.com/spf13/cobra"

// Scheduler Entity Commands
func SchedulerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduler",
		Short: "Commands related to the scheduler service",
	}
	return cmd
}

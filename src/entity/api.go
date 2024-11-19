package entity

import (
	"log"

	"github.com/spf13/cobra"
)

// API Service Commands
func APICommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Commands for API service management",
	}

	cmd.AddCommand(startApiCommand())
	cmd.AddCommand(statusApiCommand())
	cmd.AddCommand(stopApiCommand())
	return cmd
}

func startApiCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start API service",
		Run: func(cmd *cobra.Command, args []string) {
			server, _ := cmd.Flags().GetString("server")
			port, _ := cmd.Flags().GetInt("port")
			log.Printf("Starting API service on server %s at port %d\n", server, port)
		},
	}
	cmd.Flags().String("server", "", "Server container name")
	cmd.Flags().Int("port", 0, "Port for the API service")
	return cmd
}

func statusApiCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check API service status",
		Run: func(cmd *cobra.Command, args []string) {
			server, _ := cmd.Flags().GetString("server")
			log.Printf("Checking status of API service on server %s\n", server)
		},
	}
	cmd.Flags().String("server", "", "Server container name")
	return cmd
}

func stopApiCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop API service",
		Run: func(cmd *cobra.Command, args []string) {
			server, _ := cmd.Flags().GetString("server")
			log.Printf("Stopping API service on server %s\n", server)
		},
	}
	cmd.Flags().String("server", "", "Server container name")
	return cmd
}

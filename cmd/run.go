package cmd

import (
	"github.com/spf13/cobra"
	"synod/api"
	"synod/conf"
)

var apiService *api.ObjectServer

func RunCommand() *cobra.Command {
	runner := &cobra.Command{
		Use:   "run [server]",
		Short: "run a service",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := conf.Startup(); err != nil {
				return err
			}

			return nil
		},
	}

	api := &cobra.Command{
		Use:   "api",
		Short: "run api service",
		PreRun: func(cmd *cobra.Command, args []string) {
			apiService = api.NewObjectServer()
		},
		Run: func(cmd *cobra.Command, args []string) {
			apiService.Run()
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			apiService.Close()
		},
	}

	runner.AddCommand(api)

	return runner
}

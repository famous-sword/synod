package cmd

import (
	"github.com/spf13/cobra"
	"synod/api"
	"synod/conf"
	"synod/storage"
)

var apiService *api.Service

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
			apiService = api.New()
		},
		Run: func(cmd *cobra.Command, args []string) {
			apiService.Run()
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			apiService.Close()
		},
	}

	storage := &cobra.Command{
		Use:   "storage",
		Short: "run storage service",
		Run: func(cmd *cobra.Command, args []string) {
			svc := storage.New()
			svc.Run()
		},
	}

	runner.AddCommand(api, storage)

	return runner
}

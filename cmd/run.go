package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"synod/api"
	"synod/storage"
	"synod/util/logx"
	"syscall"
)

func RunCommand() *cobra.Command {
	runner := &cobra.Command{
		Use:              "run [server]",
		Short:            "run a service",
		TraverseChildren: true,
	}

	runner.AddCommand(newAPICommand(), newStorageCommand())

	return runner
}

func newAPICommand() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "start a api sever",
		Run: func(cmd *cobra.Command, args []string) {
			svc := api.New()
			quit := make(chan os.Signal)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				err := svc.Run()
				if err != nil {
					logx.Errorw("api service run error", "error", err)
				}
			}()

			<-quit

			svc.Shutdown()
		},
	}
}

func newStorageCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "storage",
		Short: "start a storage server",
		Run: func(cmd *cobra.Command, args []string) {
			svc := storage.New()
			quit := make(chan os.Signal)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				err := svc.Run()

				if err != nil {
					logx.Errorw("storage service run error", "error", err)
				}
			}()

			<-quit

			svc.Shutdown()
		},
	}
}

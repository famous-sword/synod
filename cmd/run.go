package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"synod/api"
	"synod/conf"
	"synod/storage"
	"syscall"
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
		Run: func(cmd *cobra.Command, args []string) {
			quit := make(chan os.Signal)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			svc := api.New()
			go func() {
				err := svc.Run()

				if err != nil {
					quit <- syscall.SIGINT
					log.Fatalln(err)
				}
			}()

			<-quit
			svc.Shutdown()
		},
	}

	storage := &cobra.Command{
		Use:   "storage",
		Short: "run storage service",
		Run: func(cmd *cobra.Command, args []string) {
			quit := make(chan os.Signal)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			svc := storage.New()

			go func() {
				err := svc.Run()

				if err != nil {
					quit <- syscall.SIGINT
					log.Println(err)
				}
			}()

			<-quit

			svc.Shutdown()
		},
	}

	runner.AddCommand(api, storage)

	return runner
}

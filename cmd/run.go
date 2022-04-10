package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"synod/api"
	"synod/conf"
	"synod/storage"
	"synod/util/logx"
	"syscall"
)

func RunCommand() *cobra.Command {
	runner := &cobra.Command{
		Use:   "run [server]",
		Short: "run a service",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if err = conf.Startup(); err != nil {
				return err
			}

			if err = logx.Setup(); err != nil {
				return err
			}

			if conf.Bool("app.debug") {
				gin.SetMode(gin.DebugMode)
			} else {
				gin.SetMode(gin.ReleaseMode)
			}

			return nil
		},
	}

	runner.AddCommand(APICommand(), StorageCommand())

	return runner
}

func APICommand() *cobra.Command {
	return &cobra.Command{
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
					logx.Errorw("api service run error", err)
				}
			}()

			<-quit
			svc.Shutdown()
		},
	}
}

func StorageCommand() *cobra.Command {
	return &cobra.Command{
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
					logx.Errorw("storage service run error", err)
				}
			}()

			<-quit

			svc.Shutdown()
		},
	}
}

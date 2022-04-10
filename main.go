package main

import (
	"github.com/spf13/cobra"
	"synod/cmd"
	"synod/discovery"
	"synod/util/logx"
)

func main() {
	root := &cobra.Command{
		Use:   "synod",
		Short: "Simple distributed object storage system",
	}

	root.AddCommand(cmd.RunCommand())

	if err := root.Execute(); err != nil {
		logx.Errorw("synod exec error", err)
	}

	discovery.Shutdown()
	logx.Flush()
}

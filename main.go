package main

import (
	"github.com/spf13/cobra"
	"log"
	"synod/cmd"
)

func main() {
	root := &cobra.Command{
		Use:   "synod",
		Short: "Simple distributed object storage system",
	}

	root.AddCommand(cmd.RunCommand())

	if err := root.Execute(); err != nil {
		log.Println(err)
	}
}

package cmd

import (
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"synod/conf"
	"synod/discovery"
	"synod/util/logx"
)

// use spf13/cobra，read flag always empty
// use the native flag for now
// wait for the spf13/cobra next version to try
var configFile = flag.String("c", "", "set config file")

type Synod struct {
	name        string
	version     string
	intro       string
	configFile  string
	rootCommand *cobra.Command
	subCommands []*cobra.Command
}

func NewSynod() *Synod {
	return &Synod{
		name:    "synod",
		version: "0.2-dev",
		intro:   "Simple distributed object storage system",
	}
}

func (s *Synod) Run() {
	flag.Parse()
	s.startupRootCommand()
	s.startupFeatureCommands()
	s.startupFlagParser()
	s.setupCoreComponentsOrStop()

	if err := s.rootCommand.Execute(); err != nil {
		logx.Fatalw("synod run error", "message", err)
	}
}

func (s *Synod) Shutdown() {
	discovery.Shutdown()
	logx.Flush()
}

func (s *Synod) startupRootCommand() {
	s.rootCommand = &cobra.Command{
		Use:              s.name,
		Version:          s.version,
		Short:            s.intro,
		TraverseChildren: true,
	}
}

func (s *Synod) startupFeatureCommands() {
	s.addSubCommand(RunCommand())

	s.rootCommand.AddCommand(s.subCommands...)
}

func (s *Synod) startupFlagParser() {
	s.rootCommand.PersistentFlags().StringVarP(
		&s.configFile,
		"config",
		"c",
		"",
		"set config file",
	)
}

func (s *Synod) setupCoreComponentsOrStop() {
	var err error
	if err = conf.Startup(*configFile); err != nil {
		fmt.Printf("setup config error: %v\n", err)
		os.Exit(1)
	}
	if err = logx.Startup(); err != nil {
		fmt.Printf("setup logx error: %v\n", err)
		os.Exit(1)
	}
}

func (s *Synod) addSubCommand(command *cobra.Command) {
	s.subCommands = append(s.subCommands, command)
}

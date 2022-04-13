package main

import (
	"synod/cmd"
)

func main() {
	app := cmd.NewSynod()
	app.Run()
	app.Shutdown()
}

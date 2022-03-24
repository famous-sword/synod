package main

import (
	"log"
	"os"
	"os/signal"
	"synod/api"
	"synod/conf"
	"syscall"
)

func main() {
	if err := conf.Startup(); err != nil {
		log.Fatalln(err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	app := api.NewObjectServer()

	go func() {
		app.Run()
	}()

	<-quit
	app.Close()
}

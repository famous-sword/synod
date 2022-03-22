package main

import (
	"log"
	"os"
	"os/signal"
	"synod/api"
	"synod/conf"
	"synod/discovery"
	"syscall"
)

func main() {
	if err := conf.Startup(); err != nil {
		log.Fatalln(err)
	}

	if err := discovery.Startup(); err != nil {
		log.Fatalln(err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := discovery.StartHeartbeat()
		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		err := api.Run()
		if err != nil {
			log.Println(err)
		}
	}()

	<-quit
	discovery.Close()
}

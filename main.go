package main

import (
	"log"
	"os"
	"os/signal"
	"synod/conf"
	"synod/discovery"
	"syscall"
)

func main() {
	if err := conf.Startup(); err != nil {
		log.Fatalln(err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	//go func() {
	//	err := api.Run()
	//	if err != nil {
	//		log.Println(err)
	//	}
	//}()

	publisher := discovery.NewPublisher("api", "localhost:8000")

	if err := publisher.Publish(); err != nil {
		log.Println(err)
	}

	<-quit
	log.Println(publisher.Unpublished())
}

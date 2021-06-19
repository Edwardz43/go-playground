package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/reactivex/rxgo/v2"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ch := make(chan rxgo.Item)
	go factroy(ch)
	// Create a Connectable Observable
	observable := rxgo.FromChannel(ch, rxgo.WithPublishStrategy()).Filter(func(i interface{}) bool {
		return i.(int) > 2
	})

	// Create the first Observer
	observable.DoOnNext(func(i interface{}) {
		fmt.Printf("First observer: %d\n", i)
	})

	// Create the second Observer
	observable.DoOnNext(func(i interface{}) {
		fmt.Printf("Second observer: %d\n", i)
	})

	observable.Connect()
	<-quit
	log.Println("service shutdown")
	os.Exit(1)
}

func factroy(ch chan rxgo.Item) {
	ch <- rxgo.Of(1)
	ch <- rxgo.Of(2)
	ch <- rxgo.Of(3)
	close(ch)
}

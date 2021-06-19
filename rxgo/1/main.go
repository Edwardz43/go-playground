package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/reactivex/rxgo/v2"
)

type Customer struct {
	ID             int
	Name, LastName string
	Age            int
	TaxNumber      string
}

func (c Customer) toString() string {
	return fmt.Sprintf("ID : %d, Name : %s, LastName : %s, Age : %d, TaxNumber : %s", c.ID, c.Name, c.LastName, c.Age, c.TaxNumber)
}

var (
	pool int = 4
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Create the input channel
	ch := make(chan rxgo.Item)
	done := make(chan struct{})

	//Data producer
	go producer(quit, done, ch)

	// Create an Observable
	observable := rxgo.
		FromEventSource(ch).
		Filter(func(item interface{}) bool {
			// Filter operation
			return item.(Customer).Age > 18
		}).
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			// Enrich operation
			customer := item.(Customer)
			taxNumber, _ := getTaxNumber(customer)
			customer.TaxNumber = taxNumber
			return customer, nil
		}, rxgo.WithPublishStrategy())

	// for i := range observable.Observe() {
	// 	log.Printf("observe : %v", i.V)
	// }

	observable.DoOnNext(func(i interface{}) {
		fmt.Printf("observer: \x1b[34m%v\x1b[0m\n", i.(Customer).toString())
	})

	observable.Connect()

	select {
	case <-done:
		log.Println("service shutdown")
		os.Exit(0)
	}
}

func producer(quit chan os.Signal, done chan struct{}, ch chan rxgo.Item) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			i := rxgo.Item{
				V: Customer{
					ID:        time.Now().Nanosecond(),
					Name:      fmt.Sprintf("Name%d", rand.Intn(100)),
					LastName:  fmt.Sprintf("LastName%d", rand.Intn(100)),
					Age:       rand.Intn(65),
					TaxNumber: "",
				},
				E: nil,
			}
			log.Printf("produce customer %s, age=%d", i.V.(Customer).Name, i.V.(Customer).Age)
			ch <- i
		case <-quit:
			ticker.Stop()
			close(ch)
			d := new(struct{})
			done <- *d
		}
	}
}

func getTaxNumber(c Customer) (string, error) {
	return fmt.Sprintf("TAX%s", time.Now().Format("20060102150405")), nil
}

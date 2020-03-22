package main

import (
	"log"
	"time"
)

func main() {

	num := 20000
	// sig, cancel := context.WithCancel(context.Background())
	start := time.Now()
	ch := generateNatural()
	for i := 0; i < num; i++ {
		prime := <-ch
		// fmt.Printf("%v:%v\n", i+1, prime)
		ch = primeFilter(ch, prime)
	}
	log.Printf("%d prime is %v", num, <-ch)
	log.Printf("spand %v", time.Since(start))

}

func generateNatural() chan int {
	ch := make(chan int)
	go func() {
		for i := 2; ; i++ {
			ch <- i
			// time.Sleep(time.Millisecond * 100)
		}
	}()
	return ch
}

func primeFilter(in <-chan int, prime int) chan int {
	out := make(chan int)

	go func() {
		for {
			if i := <-in; i%prime != 0 {
				out <- i
			}
		}
	}()
	return out
}

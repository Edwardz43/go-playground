package main

import (
	"sync"
	"sync/atomic"
)

var counter int32 = 0

var lock sync.Mutex

func main() {
	cas()
}

func cas() {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt32(&counter, 1)
		}()
	}

	wg.Wait()
	println(counter)
}

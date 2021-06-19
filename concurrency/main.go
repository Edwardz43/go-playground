package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

var (
	jobCount     uint64 = 0
	worker1Count uint64 = 0
	worker2Count uint64 = 0
	worker3Count uint64 = 0
	worker4Count uint64 = 0

	cpuprofile    = flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile    = flag.String("memprofile", "", "write memory profile to `file`")
	workerCounter = flag.Uint64("workercounter", 5, "worker counter")
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ch := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan string)
	start := time.Now()
	go func() {
		for i := 0; ; i++ {
			if atomic.LoadUint64(&worker1Count) < *workerCounter {
				atomic.AddUint64(&worker1Count, 1)
				go func(i int) {
					log.Printf("%s[%d]init ...%s\n", colorCyan, i, colorWhite)
					// r := fmt.Sprintf("%ds", rand.Intn(3))
					// s, _ := time.ParseDuration(r)
					time.Sleep(1 * time.Second)
					ch <- 1
					atomic.AddUint64(&worker1Count, ^uint64(0))
				}(i)
			}
		}
	}()
	go step1(ch, ch2)
	go step2(ch2, ch3)
	go step3(ch3, done)

	go func() {
		_ = <-sigs
		fmt.Println("got the stop channel")
		fmt.Printf("%d, %d, %d\n", atomic.LoadUint64(&worker1Count), atomic.LoadUint64(&worker2Count), atomic.LoadUint64(&worker3Count))
		done <- true
	}()
	<-done
	end := time.Since(start)
	log.Println(end)
	close(ch)
	close(ch2)
	close(ch3)
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

func step1(ch <-chan int, ch2 chan<- int) {
	for i := range ch {
		if atomic.LoadUint64(&worker2Count) < *workerCounter {
			// fmt.Printf("worker2[%d]\n", atomic.LoadUint64(&worker2Count))
			atomic.AddUint64(&worker2Count, 1)
			go func() {
				log.Printf("%sstep1 receive data%s\n", colorRed, colorWhite)
				// r := fmt.Sprintf("%ds", rand.Intn(5))
				// s, _ := time.ParseDuration(r)
				// time.Sleep(s)
				time.Sleep(3 * time.Second)
				ch2 <- i * time.Now().Nanosecond()
				atomic.AddUint64(&worker2Count, ^uint64(0))
			}()
		}
	}
}

func step2(ch2 <-chan int, ch3 chan<- string) {
	for s := range ch2 {
		if atomic.LoadUint64(&worker3Count) < *workerCounter {
			// fmt.Printf("worker3[%d]\n", worker3Count)
			atomic.AddUint64(&worker3Count, 1)
			go func(s string) {
				log.Printf("%sstep2 receive data[%s]%s\n", colorPurple, s, colorWhite)
				// r := fmt.Sprintf("%ds", rand.Intn(3))
				// d, _ := time.ParseDuration(r)
				time.Sleep(2 * time.Second)
				ch3 <- s
				atomic.AddUint64(&worker3Count, ^uint64(0))
			}(strconv.Itoa(s))
		}
	}
}

func step3(ch3 <-chan string, done chan bool) {
	if atomic.LoadUint64(&worker4Count) < *workerCounter {

		atomic.AddUint64(&worker4Count, 1)
		for _ = range ch3 {
			// fmt.Printf("result = %s\n", r)
			atomic.AddUint64(&jobCount, 1)
			log.Printf("result[%d]\n", atomic.LoadUint64(&jobCount))
			if atomic.LoadUint64(&jobCount) >= 100 {
				done <- true
				return
			}
			// wg.Done()
			atomic.AddUint64(&worker4Count, ^uint64(0))
		}
	}
}

package main

import (
	"flag"
	"log"
	"os/exec"
)

var url = flag.String("u", "localhost", "url")
var args = flag.String("a", "-type=ns", "arguments")

func main() {
	flag.Parse()
	e := exec.Command("nslookup", *args, *url)

	var msg []byte
	var err error

	if msg, err = e.Output(); err != nil {
		log.Fatal("2 :", err)
	}

	log.Println(string(msg))
}

package main

import (
	"log"
	"net/http"
)

var (
	errCount          uint64 = 0
	normalCount       uint64 = 0
	totalResponseTime int64  = 0
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("Hello World")
	})

	if err :=http.ListenAndServe(":8091", nil); err != nil{
		log.Panicln(err.Error())
	}
}

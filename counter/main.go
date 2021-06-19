package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"os"
)

type DB struct {
	*redis.Client
}

var (
	db       DB
	hostname string
)

func (db *DB) get() (int64, bool) {
	s := db.Get("mycounter")
	c, err := s.Int64()
	if err != nil {
		return 0, false
	}
	return c, true
}

func (db *DB) incr() (int64, bool) {
	s := db.Incr("mycounter")
	c, err := s.Result()
	if err != nil {
		return 0, false
	}
	return c, true
}

//

func main() {
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/incr", incrHandler)

	log.Println("service runs at http://localhost:80")
	_ = http.ListenAndServe(":80", nil)
}

func getHandler(writer http.ResponseWriter, request *http.Request) {
	c, ok := db.get()
	if !ok {
		log.Println("get error")
		_, _ = writer.Write([]byte("get error"))
		return
	}
	msg := fmt.Sprintf("[%s] %d", hostname, c)
	_, _ = writer.Write([]byte(msg))
}

func incrHandler(writer http.ResponseWriter, request *http.Request) {
	c, ok := db.incr()
	if !ok {
		log.Println("incr error")
		_, _ = writer.Write([]byte("incr error"))
		return
	}
	msg := fmt.Sprintf("[%s] %d", hostname, c)
	_, _ = writer.Write([]byte(msg))
}

func init() {
	hostname, _ = os.Hostname()
	conn := redis.NewClient(&redis.Options{
		Addr: "redis-1594901280-master.default.svc.cluster.local:6379",
		// Addr:     "localhost:42529",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	s := conn.Ping()
	pong, ok := s.Result()
	if ok != nil {
		log.Panic(ok)
	}
	log.Println(pong)

	db = DB{conn}
}

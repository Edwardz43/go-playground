package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {

	flag.Parse()
	server := newServer()
	go server.run()

	go func() {
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			ServeWs(server, w, r)
		})
		log.Printf("serve at http://localhost%s", *addr)
		err := http.ListenAndServe(*addr, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}

	}()
	//http.HandleFunc("/", serveHome)
	time.Sleep(time.Second)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c := Conn()
	defer c.Close()

	process(c, interrupt)

}

func process(c *websocket.Conn, interrupt chan os.Signal) {
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				print(fmt.Sprintf("client read error : %s", err))
				return
			}
			print(message)
		}
	}()

	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')

			text := strings.Replace(input, "\n", "", 1)
			err := c.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				print(fmt.Sprintf("write message => %s", err))
				break
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return

		case <-interrupt:
			print("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				print(fmt.Sprintf("write close: %s", err))
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func print(msg interface{}) {
	fmt.Printf("\x1b[%dmclient : %s\x1b[0m\n", 33, msg)
}

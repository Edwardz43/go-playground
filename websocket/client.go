package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	id     uint64
	server *Server         // The websocket connection.
	conn   *websocket.Conn // Buffered channel of outbound messages.
	send   chan []byte
}

func sPrint(msg interface{}) {
	fmt.Printf("\x1b[%dmclient : %s\x1b[0m\n", 32, msg)
}

// read pumps messages from the websocket connection to the hub.
//
// The application runs read in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) read() {
	defer func() {
		sPrint("Client connection closed")
		c.server.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)

	// c.conn.WritePreparedMessage()()
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		// log.Println("server pong")
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				sPrint(fmt.Sprintf("error: %v", err))
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.server.broadcast <- message
	}
}

// write pumps messages from the hub to the websocket connection.
//
// A goroutine running write is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	timerTicker := time.NewTicker(time.Second * 5)
	defer func() {
		ticker.Stop()
		timerTicker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := writeMsg(c.conn, message, c.send); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// log.Println("server ping")
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-timerTicker.C:
			datetime := time.Now().UTC().Format("20060102 15:04:05")

			if err := writeMsg(c.conn, []byte(datetime), c.send); err != nil {
				return
			}
		}
	}
}

func writeMsg(c *websocket.Conn, msg []byte, ch chan []byte) error {

	w, err := c.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	w.Write([]byte(msg))

	// Add queued chat messages to the current websocket message.
	n := len(ch)
	for i := 0; i < n; i++ {
		log.Printf("queued chat message : [%s]", <-ch)
		w.Write(newline)
		w.Write(<-ch)
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

// Conn creates and returns a websocket connection.
func Conn() *websocket.Conn {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	print(fmt.Sprintf("connecting to %s", u.String()))

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return c
}

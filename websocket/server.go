package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Server represents a websocket server.
type Server struct {
	clientID uint64

	// Registered clients.
	clients map[uint64]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func newServer() *Server {
	return &Server{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uint64]*Client),
	}
}

func (s *Server) run() {
	for {
		select {
		case client := <-s.register:
			s.clients[client.id] = client
			log.Printf("client register id :[%d]\n", client.id)
		case client := <-s.unregister:
			if _, ok := s.clients[client.id]; ok {
				delete(s.clients, client.id)
				close(client.send)
				sPrint(fmt.Sprintf("client unregister id :[%d]", client.id))
			}
		case message := <-s.broadcast:
			for id, client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, id)
				}
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(server *Server, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		sPrint(fmt.Sprintf("server:%s", err))
		return
	}
	client := &Client{
		server: server,
		conn:   conn,
		send:   make(chan []byte, 256),
	}
	client.server.clientID++

	client.id = client.server.clientID

	client.server.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.write()
	go client.read()
}

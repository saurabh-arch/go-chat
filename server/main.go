package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

// Message struct to hold incoming/outgoing messages
type Message struct {
	Text string    `json:"text"`
	Time time.Time `json:"time"`
}

type hub struct {
	clients          map[string]*websocket.Conn
	addClientChan    chan *websocket.Conn
	removeClientChan chan *websocket.Conn
	broadcastChan    chan Message
}

var (
	port = flag.String("port", "9000", "port used for ws connection")
)

/*
The server function routes its sole path (“/”) to handler, so create that function. I pass it the websocket.Conn and hub. Every client that connects will call this handler.
*/
func server(port string) error {
	h := newHub()
	mux := http.NewServeMux()
	mux.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		handler(ws, h)
	}))
	srvr := http.Server{Addr: ":" + port, Handler: mux}
	return srvr.ListenAndServe()
}

func newHub() *hub {
	return &hub{
		clients:          make(map[string]*websocket.Conn),
		addClientChan:    make(chan *websocket.Conn),
		removeClientChan: make(chan *websocket.Conn),
		broadcastChan:    make(chan Message),
	}
}

/*
handler function calls hub.run and then listens for messages from that client in a for loop.
*/
func handler(ws *websocket.Conn, h *hub) {
	go h.run()
	h.addClientChan <- ws
	for {
		var m Message
		err := websocket.JSON.Receive(ws, &m)
		if err != nil {
			h.broadcastChan <- Message{err.Error(), time.Now()}
			h.removeClient(ws)
			return
		}
		h.broadcastChan <- m
	}
}

// This listens to all the hub’s various channels, via Go’s idiomatic for-select, calling the appropriate method for each incoming channel.
func (h *hub) run() {
	for {
		select {
		case conn := <-h.addClientChan:
			h.addClient(conn)
		case conn := <-h.removeClientChan:
			h.removeClient(conn)
		case m := <-h.broadcastChan:
			h.broadcastMessage(m)
		}
	}
}

/*
Create the three methods called in run. The addClient and removeClient methods add & remove a conn from the pool (called clients in the struct) respectively. The broadcastMessage method calls websocket’s JSON.Send() method for each client in the pool.
*/

func (h *hub) removeClient(conn *websocket.Conn) {
	delete(h.clients, conn.LocalAddr().String())
}
func (h *hub) addClient(conn *websocket.Conn) {
	h.clients[conn.RemoteAddr().String()] = conn
}

func (h *hub) broadcastMessage(m Message) {
	for _, conn := range h.clients {
		err := websocket.JSON.Send(conn, m)
		if err != nil {
			fmt.Println("Error broadcasting message: ", err)
			return
		}
	}
}

/*
create main(), in which we’ll call flag.Parse() to get our flag-ified port value and call server(), running the whole shebang.
*/
func main() {
	flag.Parse()
	log.Fatal(server(*port))
}

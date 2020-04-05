package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

// Message struct to hold message related data
type Message struct {
	Name string    `json:"name"`
	Text string    `json:"text"`
	Time time.Time `json:"time"`
}

// Client struct to hold client related data
type Client struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

var (
	port   = flag.String("port", "9000", "port used for ws connection")
	name   = flag.String("name", "", "name of client") // default value - empty string
	layout = "Mon, 02 Jan 15:04:05 IST"
)

func connect() (*websocket.Conn, error) {
	return websocket.Dial(fmt.Sprintf("ws://localhost:%s", *port), "", mockedIP())
}

// As we’re planning to run all this (server and clients) locally for the demo, we need a way to differentiate the clients and can’t use localhost as the 3rd parameter (the origin) to websocket.Dial(), since every client will be localhost. I created and, above, called a mockedIP() function that just returns a faux IP as a string, for the sake of this exercise.
func mockedIP() string {
	var arr [4]int
	for i := 0; i < 4; i++ {
		rand.Seed(time.Now().UnixNano())
		arr[i] = rand.Intn(256)
	}
	return fmt.Sprintf("http://%d.%d.%d.%d", arr[0], arr[1], arr[2], arr[3])
}

// Register method to register client to server
func (c *Client) Register(serverAddr string, clientAddr string) (bool, error) {

	return true, nil
}

func getLocalIPAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return err.Error()
	}

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return fmt.Sprintf("%s", localAddr)
}

func main() {
	flag.Parse()
	if strings.TrimSpace(*name) == "" { // if no value is passed as flag parameter for name
		fmt.Println("Enter name")
		fmt.Scanf("%s", name)
	}
	if strings.TrimSpace(*name) == "" { // if name still equals empty string after user input also, assign ip address of user
		*name = getLocalIPAddress()
	}
	// connect
	ws, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	// receive
	var m Message
	go func() {
		for {
			err := websocket.JSON.Receive(ws, &m)
			if err != nil {
				fmt.Println("Error receiving message: ", err.Error())
				break
			}
			fmt.Println(m.Name, ":", m.Time.Format(layout), m.Text)
		}
	}()
	// send
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		m := Message{
			Name: *name,
			Text: text,
			Time: time.Now(),
		}
		err = websocket.JSON.Send(ws, m)
		if err != nil {
			fmt.Println("Error sending message: ", err.Error())
			break
		}
	}
}

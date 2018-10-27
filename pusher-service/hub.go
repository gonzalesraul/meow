package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

//Hub is a structure to manage all the clients
type Hub struct {
	clients    []*Client
	nextID     int
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

func newHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		nextID:     0,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}

func (hub *Hub) run() Hub {
	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

func (hub *Hub) broadcast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, c := range hub.clients {
		if c != ignore {
			c.outbound <- data
		}
	}
}

func (hub *Hub) send(message interface{}, client *Client) {
	data, _ := json.Marshal(message)
	client.outbound <- data
}

func (hub *Hub) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not upgrade", http.StatusInternalServerError)
		return
	}
	client := newClient(hub, socket)
	hub.register <- client

	go client.write()
}

func (hub *Hub) onConnect(client *Client) {
	log.Println("client connected: ", client.socket.RemoteAddr())

	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	client.id = hub.nextID
	hub.nextID++
	hub.clients = append(hub.clients, client)
}

func (hub *Hub) onDisconnect(client *Client) {
	log.Println("client disconnected: ", client.socket.RemoteAddr())

	client.close()
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	//Iterate through all the clients and remove it by its ID.
	for i, cur := range hub.clients {
		if cur.id == client.id {
			hub.clients = append(hub.clients[:i], hub.clients[i+1:]...)
			break
		}
	}
}

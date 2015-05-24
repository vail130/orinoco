package sieve

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var ActiveClients = make(map[ClientConn]int)
var ActiveClientsRWMutex sync.RWMutex

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  net.Addr
}

func addClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	ActiveClients[cc] = 0
	ActiveClientsRWMutex.Unlock()
}

func deleteClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	delete(ActiveClients, cc)
	ActiveClientsRWMutex.Unlock()
}

func broadcastMessage(messageType int, message []byte) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()

	for client, _ := range ActiveClients {
		if err := client.websocket.WriteMessage(messageType, message); err != nil {
			return
		}
	}
}

func Sieve(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/events/{event}", EventHandler).Methods("POST")
	r.HandleFunc("/subscribe", SubscribeHandler)
	http.Handle("/", r)

	var buffer bytes.Buffer
	buffer.WriteString(":")
	buffer.WriteString(port)
	port = buffer.String()

	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func EventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := vars["event"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	
	event_data := append([]byte(event), []byte(": ")...)
	event_data = append(event_data, data...)

	broadcastMessage(websocket.TextMessage, event_data)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := ws.RemoteAddr()
	sockCli := ClientConn{ws, client}
	addClient(sockCli)
}

package sieve

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const websocketBufferSize = 1024

var upgrader = websocket.Upgrader{
	ReadBufferSize:  websocketBufferSize,
	WriteBufferSize: websocketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		// TODO Actually check origin
		return true
	},
}

var sieveBoundary []byte
var ActiveClients = make(map[net.Addr]*websocket.Conn)
var ActiveClientsRWMutex sync.RWMutex

func addClient(ws *websocket.Conn) {
	ActiveClientsRWMutex.Lock()
	ActiveClients[ws.RemoteAddr()] = ws
	ActiveClientsRWMutex.Unlock()
}

func deleteClient(ws *websocket.Conn) {
	ActiveClientsRWMutex.Lock()
	delete(ActiveClients, ws.RemoteAddr())
	ActiveClientsRWMutex.Unlock()
}

func broadcastOverWebsocket(messageType int, message []byte) {
	ActiveClientsRWMutex.Lock()
	defer ActiveClientsRWMutex.Unlock()

	message = append(message, sieveBoundary...)

	for _, ws := range ActiveClients {
		if err := ws.WriteMessage(messageType, message); err != nil {
			log.Fatalln(err)
			deleteClient(ws)
		}
	}
}

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
	addClient(ws)
}

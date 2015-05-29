package sieve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

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
var ActiveClients = make(map[string](map[net.Addr]*websocket.Conn))
var ActiveClientsRWMutex sync.RWMutex

func addClient(clientType string, ws *websocket.Conn) {
	ActiveClientsRWMutex.Lock()
	ActiveClients[clientType][ws.RemoteAddr()] = ws
	ActiveClientsRWMutex.Unlock()
}

func deleteClient(clientType string, ws *websocket.Conn) {
	ActiveClientsRWMutex.Lock()
	delete(ActiveClients[clientType], ws.RemoteAddr())
	ActiveClientsRWMutex.Unlock()
}

func readMessage(boundary string) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()

	fmt.Println("READING MESSAGES")

	for _, ws := range ActiveClients["publish"] {
		boundaryBytes := []byte(boundary)
		var fullMessage []byte

		for {
			_, partialMessage, err := ws.ReadMessage()
			if err != nil {
				log.Println(err)
				deleteClient("publish", ws)
			}

			fullMessage = append(fullMessage, partialMessage...)

			fmt.Println(string(fullMessage))

			if bytes.Index(fullMessage, boundaryBytes) > -1 {
				messagePieces := bytes.SplitN(fullMessage, boundaryBytes, 1)
				fullMessage = messagePieces[0][:len(messagePieces[0])-len(boundaryBytes)]
				break
			}
		}

		var event Event
		err := json.Unmarshal(fullMessage, &event)
		timestamp, err := time.Parse(time.RFC3339, event.Timestamp)
		if err != nil {
			timestamp = time.Now()
		}
		processEvent(event.Event, timestamp, event.Data)
	}
}

func broadcastMessage(messageType int, message []byte) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()

	message = append(message, sieveBoundary...)

	for _, ws := range ActiveClients["subscribe"] {
		if err := ws.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			deleteClient("subscribe", ws)
		}
	}
}

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	addClient("subscribe", ws)
}

func PublishHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	addClient("publish", ws)
}

package sieve

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func readMessage(ws *websocket.Conn) error {
	_, message, err := ws.ReadMessage()
	message = message[:len(message)-len(sieveBoundary)]

	var event Event
	err = json.Unmarshal(message, &event)
	timestamp, err := time.Parse(time.RFC3339, event.Timestamp)
	if err != nil {
		timestamp = time.Now()
	}
	return processEvent(event.Event, timestamp, event.Data)

}

func PublishHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}

	for {
		err := readMessage(ws)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

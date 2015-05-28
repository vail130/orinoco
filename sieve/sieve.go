package sieve

import (
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/vail130/orinoco/stringutils"
)

var isTestEnv bool

func Sieve(port string, boundary string) {
	sieveBoundary = []byte(boundary)
	isTestEnv = stringutils.StringToBool(os.Getenv("TEST"))
	
	ActiveClients["subscribe"] = make(map[net.Addr]*websocket.Conn)
	ActiveClients["publish"] = make(map[net.Addr]*websocket.Conn)

	r := mux.NewRouter()
	r.HandleFunc("/events", GetAllEventsHandler).Methods("GET")
	r.HandleFunc("/events", DeleteAllEventsHandler).Methods("DELETE")
	r.HandleFunc("/events/{event}", GetEventHandler).Methods("GET")
	r.HandleFunc("/events/{event}", PostEventHandler).Methods("POST")
	r.HandleFunc("/subscribe", SubscribeHandler)
	r.HandleFunc("/publish", PublishHandler)
	http.Handle("/", r)

	port = stringutils.Concat(":", port)

	err := http.ListenAndServe(port, nil)
	if err != nil && err != syscall.EPIPE {
		panic("ListenAndServe: " + err.Error())
	}
}

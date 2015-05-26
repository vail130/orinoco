package sieve

import (
	"fmt"
	"net/http"
	"syscall"

	"github.com/gorilla/mux"
	
	"../stringutils"
)

func Sieve(port string, boundary string) {
	sieveBoundary = []byte(boundary)
	
	r := mux.NewRouter()
	r.HandleFunc("/events", GetAllEventsHandler).Methods("GET")
	r.HandleFunc("/events/{event}", GetEventHandler).Methods("GET")
	r.HandleFunc("/events/{event}", PostEventHandler).Methods("POST")
	r.HandleFunc("/subscribe", SubscribeHandler)
	http.Handle("/", r)
	
	fmt.Println("Running Sieve server on port", port)
	
	port = stringutils.Concat(":", port)

	err := http.ListenAndServe(port, nil)
	if err != nil && err != syscall.EPIPE {
		panic("ListenAndServe: " + err.Error())
	}
}

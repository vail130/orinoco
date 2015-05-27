package sieve

import (
	"net/http"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/vail130/orinoco/stringutils"
)

func Sieve(port string, boundary string) {
	sieveBoundary = []byte(boundary)

	r := mux.NewRouter()
	r.HandleFunc("/events", GetAllEventsHandler).Methods("GET")
	r.HandleFunc("/events", DeleteAllEventsHandler).Methods("DELETE")
	r.HandleFunc("/events/{event}", GetEventHandler).Methods("GET")
	r.HandleFunc("/events/{event}", PostEventHandler).Methods("POST")
	r.HandleFunc("/subscribe", SubscribeHandler)
	http.Handle("/", r)

	port = stringutils.Concat(":", port)

	err := http.ListenAndServe(port, nil)
	if err != nil && err != syscall.EPIPE {
		panic("ListenAndServe: " + err.Error())
	}
}

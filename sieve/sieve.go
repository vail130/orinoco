package sieve

import (
	"net/http"

	"github.com/gorilla/mux"
	
	"../stringutils"
)

func Sieve(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/events", GetAllEventsHandler).Methods("GET")
	r.HandleFunc("/events/{event}", GetEventHandler).Methods("GET")
	r.HandleFunc("/events/{event}", PostEventHandler).Methods("POST")
	r.HandleFunc("/subscribe", SubscribeHandler)
	http.Handle("/", r)
	
	port = stringutils.Concat(":", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}



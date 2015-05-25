package sieve

import (
	"bytes"
	"net/http"

	"github.com/gorilla/mux"
)

func Sieve(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/events/{event}", GetEventHandler).Methods("GET")
	r.HandleFunc("/events/{event}", PostEventHandler).Methods("POST")
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



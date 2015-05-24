package sieve

import (
	"fmt"
	"log"
	"net/http"
	
	"github.com/gorilla/mux"
)

func Sieve() {
	fmt.Println("Sieve")

	r := mux.NewRouter()
	r.HandleFunc("/events/{event}", EventHandler).Methods("POST")
	r.HandleFunc("/events/{event}/", EventHandler).Methods("POST")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

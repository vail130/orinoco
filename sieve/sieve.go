package sieve

import (
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/vail130/orinoco/stringutils"
)

type SieveConfig struct {
	IsTestEnv bool
	EventChannel chan *Event
}

var sieveConfig SieveConfig

func Sieve(port string, eventChannel chan *Event) {
	sieveConfig = SieveConfig{
		stringutils.StringToBool(os.Getenv("TEST")),
		eventChannel,
	}

	r := mux.NewRouter()
	s := r.PathPrefix("/streams").Subrouter()
	
	s.HandleFunc("/", GetAllStreamsHandler).Methods("GET")
	s.HandleFunc("/", DeleteAllStreamsHandler).Methods("DELETE")
	s.HandleFunc("/{stream}", GetStreamHandler).Methods("GET")
	s.HandleFunc("/{stream}", PostStreamHandler).Methods("POST")
	
	http.Handle("/", r)

	port = stringutils.Concat(":", port)

	err := http.ListenAndServe(port, nil)
	if err != nil && err != syscall.EPIPE {
		log.Fatalln(err)
	}
}

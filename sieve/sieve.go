package sieve

import (
	"net/http"
	"os"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/vail130/orinoco/stringutils"
)

var isTestEnv bool
var configuredStreams [](map[string]string)

func Sieve(port string, streams [](map[string]string)) {
	isTestEnv = stringutils.StringToBool(os.Getenv("TEST"))
	configuredStreams = streams

	r := mux.NewRouter()
	r.HandleFunc("/streams", GetAllStreamsHandler).Methods("GET")
	r.HandleFunc("/streams", DeleteAllStreamsHandler).Methods("DELETE")
	r.HandleFunc("/streams/{stream}", GetStreamHandler).Methods("GET")
	r.HandleFunc("/streams/{stream}", PostStreamHandler).Methods("POST")
	//	r.HandleFunc("/subscribe", SubscribeHandler)
	http.Handle("/", r)

	port = stringutils.Concat(":", port)

	err := http.ListenAndServe(port, nil)
	if err != nil && err != syscall.EPIPE {
		panic("ListenAndServe: " + err.Error())
	}
}

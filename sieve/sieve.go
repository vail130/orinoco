package sieve

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

//		            event       hour     minute   second  count
var eventMap = make(map[string](map[int](map[int](map[int]int))))

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

func GetEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := vars["event"]
	
	eventCount := 0
	
	if _, ok := eventMap[event]; ok == true {
		for _, minuteMap := range eventMap[event] {
			for _, secondMap := range minuteMap {
				for _, count := range secondMap {
					eventCount += count
				}
			}
		}
	}
	
	w.Write([]byte(strconv.Itoa(eventCount)))
}

func PostEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := vars["event"]
	
	t := time.Now()
	hourMap, ok := eventMap[event]
	if ok == false {
		eventMap[event] = make(map[int](map[int](map[int]int)))
		hourMap = eventMap[event]
	}
	
	minuteMap, ok := hourMap[t.Hour()]
	if ok == false {
		hourMap[t.Hour()] = make(map[int](map[int]int))
		minuteMap = hourMap[t.Hour()]
	}
	
	secondMap, ok := minuteMap[t.Minute()]
	if ok == false {
		minuteMap[t.Minute()] = make(map[int]int)
		secondMap = minuteMap[t.Minute()]
	}
	
	_, ok = secondMap[t.Second()]
	if ok == false {
		secondMap[t.Second()] = 0
	}
	secondMap[t.Second()] += 1

	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		data = make([]byte, 0)
	}

	var buffer bytes.Buffer
	buffer.WriteString(t.Format(time.RFC3339))
	buffer.WriteString(" ")
	buffer.WriteString(event)
	buffer.WriteString(" ")
	buffer.WriteString(string(data))
	event_data := buffer.String()

	broadcastMessage(websocket.TextMessage, []byte(event_data))
}

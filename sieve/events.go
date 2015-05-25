package sieve

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type EventSummary struct {
	Event                    string  `json:"event"`
	Timestamp                string  `json:"timestamp"`
	SecondToDate             int     `json:"second_to_date"`
	MinuteToDate             int     `json:"minute_to_date"`
	HourToDate               int     `json:"hour_to_date"`
	ProjectedThisMinute      float32 `json:"projected_this_minute"`
	ProjectedThisHour        float32 `json:"projected_this_hour"`
	TrailingAveragePerSecond float32 `json:"trailing_average_per_second"`
	TrailingAveragePerMinute float32 `json:"trailing_average_per_minute"`
	TrailingAveragePerHour   float32 `json:"trailing_average_per_hour"`
}

type Event struct {
	Event     string `json:"event"`
	Timestamp string `json:"timestamp"`
	Data      []byte `json:"data"`
}

// Reference time for formats: Mon Jan 2 15:04:05 -0700 MST 2006
const dateKeyFormat = "2006-01-02"
const hourKeyFormat = "2006-01-02-15"
const minuteKeyFormat = "2006-01-02-15-04"
const secondKeyFormat = "2006-01-02-15-04-05"

var dateMap = make(map[string](map[string](map[string]int)))
var dateKeyMap = make(map[string]time.Time)

func trackEventForTime(event string, t time.Time) {
	eventMap := getEventMapForTime(t)

	dateMap, ok := eventMap[event]
	if ok == false {
		eventMap[event] = make(map[string]int)
		dateMap = eventMap[event]
	}

	timeKeys := []string{
		t.Format(hourKeyFormat),
		t.Format(minuteKeyFormat),
		t.Format(secondKeyFormat),
	}

	for i := 0; i < len(timeKeys); i++ {
		if _, ok := dateMap[timeKeys[i]]; ok == false {
			dateMap[timeKeys[i]] = 0
		}
		dateMap[timeKeys[i]] += 1
	}
}

func PostEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := vars["event"]

	t := time.Now()
	trackEventForTime(event, t)

	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		data = make([]byte, 0)
	}

	eventData := Event{
		event,
		t.Format(time.RFC3339),
		data,
	}

	jsonData, err := json.Marshal(eventData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	broadcastMessage(websocket.TextMessage, jsonData)

	w.WriteHeader(http.StatusCreated)
}

func deleteObsoleteDateKeysForTime(t time.Time) {
	if t.Hour() >= 3 {
		for dateKey, tt := range dateKeyMap {
			if tt.Year() < t.Year() ||
				(tt.Year() == t.Year() && tt.Month() < t.Month()) ||
				(tt.Year() == t.Year() && tt.Month() == t.Month() && tt.Day() < t.Day()) {
				delete(dateKeyMap, dateKey)
			}
		}
	}
}

func getEventMapForTime(t time.Time) map[string](map[string]int) {
	deleteObsoleteDateKeysForTime(t)

	dateKey := t.Format(dateKeyFormat)
	eventMap, ok := dateMap[dateKey]
	if ok == false {
		dateMap[dateKey] = make(map[string](map[string]int))
		eventMap = dateMap[dateKey]
		dateKeyMap[dateKey] = t
	}

	return eventMap
}

func getEventSummary(now time.Time, eventMap map[string](map[string]int), event string) *EventSummary {
	timeUnits := map[string]time.Duration{
		"hour":   time.Hour,
		"minute": time.Minute,
		"second": time.Second,
	}

	keyFormats := map[string]string{
		"hour":   hourKeyFormat,
		"minute": minuteKeyFormat,
		"second": secondKeyFormat,
	}

	valuesToDate := map[string]int{
		"hour":   0,
		"minute": 0,
		"second": 0,
	}

	trailingCounts := map[string]int{
		"hour":   0,
		"minute": 0,
		"second": 0,
	}

	var timeKey string

	if timeMap, eventExists := eventMap[event]; eventExists {
		for period, _ := range trailingCounts {
			if periodToDate, timeKeyExists := timeMap[now.Format(keyFormats[period])]; timeKeyExists {
				valuesToDate[period] = periodToDate
			}

			onePeriodAgo := now.Add(-1 * timeUnits[period])
			timeKey = onePeriodAgo.Format(keyFormats[period])
			if timeValue, timeKeyExists := timeMap[timeKey]; timeKeyExists {
				trailingCounts[period] += timeValue
			}

			twoPeriodsAgo := now.Add(-2 * timeUnits[period])
			timeKey = twoPeriodsAgo.Format(keyFormats[period])
			if timeValue, timeKeyExists := timeMap[timeKey]; timeKeyExists {
				trailingCounts[period] += timeValue
			}

			threePeriodsAgo := now.Add(-3 * timeUnits[period])
			timeKey = threePeriodsAgo.Format(keyFormats[period])
			if timeValue, timeKeyExists := timeMap[timeKey]; timeKeyExists {
				trailingCounts[period] += timeValue
			}
		}
	}

	projectedThisMinute := float32(valuesToDate["minute"]) / float32(now.Second()+1) * float32(60)
	projectedThisHour := float32(valuesToDate["hour"]) / float32(now.Minute()+1) * float32(60)

	return &EventSummary{
		event,
		now.Format(time.RFC3339),
		valuesToDate["second"],
		valuesToDate["minute"],
		valuesToDate["hour"],
		projectedThisMinute,
		projectedThisHour,
		float32(trailingCounts["second"]) / float32(3.0),
		float32(trailingCounts["minute"]) / float32(3.0),
		float32(trailingCounts["hour"]) / float32(3.0),
	}
}

func GetEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := vars["event"]

	now := time.Now()
	eventMap := getEventMapForTime(now)
	eventResponse := *getEventSummary(now, eventMap, event)

	jsonData, err := json.Marshal(eventResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GetAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	eventMap := getEventMapForTime(now)

	var eventSummaries []EventSummary

	for event, _ := range eventMap {
		eventSummaries = append(eventSummaries, *getEventSummary(now, eventMap, event))
	}

	jsonData, err := json.Marshal(eventSummaries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

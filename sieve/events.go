package sieve

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
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
	ChangePerSecond          int     `json:"change_per_second"`
	ChangePerMinute          int     `json:"change_per_minute"`
	ChangePerHour            int     `json:"change_per_hour"`
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

func PostEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := vars["event"]

	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		data = make([]byte, 0)
	}

	t := getTimestampForRequestData(data)
	trackEventForTime(event, t)

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

func getTimestampForRequestData(data []byte) time.Time {
	if !isTestEnv {
		return time.Now()
	}
	
	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return time.Now()
	}
	
	if timestampString, ok := f.(map[string]interface{})["timestamp"]; ok {
		if timestamp, err := time.Parse(time.RFC3339, timestampString.(string)); err == nil {
			return timestamp
		}
	}
	
	return time.Now()
}

func trackEventForTime(event string, t time.Time) {
	eventMap := getEventMapForTime(t)
	deleteObsoleteDateKeysForTime(t)

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

func getEventMapForTime(t time.Time) map[string](map[string]int) {
	dateKey := t.Format(dateKeyFormat)
	eventMap, ok := dateMap[dateKey]
	if ok == false {
		dateMap[dateKey] = make(map[string](map[string]int))
		eventMap = dateMap[dateKey]
		dateKeyMap[dateKey] = t
	}

	return eventMap
}

func deleteObsoleteDateKeysForTime(t time.Time) {
	// Delete date keys more than 24 hours old
	startTime := t.AddDate(0, 0, -1)
	for dateKey, tt := range dateKeyMap {
		if tt.Year() < startTime.Year() ||
			(tt.Year() == startTime.Year() && tt.Month() < startTime.Month()) ||
			(tt.Year() == startTime.Year() && tt.Month() == startTime.Month() && tt.Day() < startTime.Day()) {
			delete(dateKeyMap, dateKey)
		}
	}
}

func GetEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := vars["event"]

	now := getTimestampForSummaryRequest(r.URL.Query())
	eventMap := getEventMapForTime(now)
	eventSummary := *getEventSummary(now, eventMap, event)

	jsonData, err := json.Marshal(eventSummary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GetAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	now := getTimestampForSummaryRequest(r.URL.Query())
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

func getTimestampForSummaryRequest(queryValues url.Values) time.Time {
	if !isTestEnv {
		return time.Now()
	}
	
	if timestampString, ok := queryValues["timestamp"]; ok {
		if timestamp, err := time.Parse(time.RFC3339, timestampString[0]); err == nil {
			return timestamp
		}
	}
	
	return time.Now()
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

	changePerPeriod := map[string]int{
		"hour":   0,
		"minute": 0,
		"second": 0,
	}

	var timeKey1 string
	var timeKey2 string
	var timeKey3 string

	if timeMap, eventExists := eventMap[event]; eventExists {
		for period, _ := range trailingCounts {
			if periodToDate, timeKeyExists := timeMap[now.Format(keyFormats[period])]; timeKeyExists {
				valuesToDate[period] = periodToDate
			}

			onePeriodAgo := now.Add(-1 * timeUnits[period])
			timeKey1 = onePeriodAgo.Format(keyFormats[period])
			timeValue1, timeKey1Exists := timeMap[timeKey1]
			if timeKey1Exists {
				trailingCounts[period] += timeValue1
			}

			twoPeriodsAgo := now.Add(-2 * timeUnits[period])
			timeKey2 = twoPeriodsAgo.Format(keyFormats[period])
			timeValue2, timeKey2Exists := timeMap[timeKey2]
			if timeKey2Exists {
				trailingCounts[period] += timeValue2
			}

			threePeriodsAgo := now.Add(-3 * timeUnits[period])
			timeKey3 = threePeriodsAgo.Format(keyFormats[period])
			timeValue3, timeKey3Exists := timeMap[timeKey3]
			if timeKey3Exists {
				trailingCounts[period] += timeValue3
			}

			if timeKey1Exists {
				changePerPeriod[period] = timeValue1
				if timeKey2Exists {
					changePerPeriod[period] = timeValue1 - timeValue2
				}
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
		changePerPeriod["second"],
		changePerPeriod["minute"],
		changePerPeriod["hour"],
	}
}

func DeleteAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	dateMap = make(map[string](map[string](map[string]int)))
	dateKeyMap = make(map[string]time.Time)
	w.WriteHeader(http.StatusNoContent)
}

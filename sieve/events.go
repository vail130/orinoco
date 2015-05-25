package sieve

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type EventResponse struct {
	SecondToDate             int     `json:"second_to_date"`
	MinuteToDate             int     `json:"minute_to_date"`
	HourToDate               int     `json:"hour_to_date"`
	TrailingAveragePerSecond float32 `json:"trailing_average_per_second"`
	TrailingAveragePerMinute float32 `json:"trailing_average_per_minute"`
	TrailingAveragePerHour   float32 `json:"trailing_average_per_hour"`
}

// Mon Jan 2 15:04:05 -0700 MST 2006
const hourKeyFormat = "2006-01-02-15"
const minuteKeyFormat = "2006-01-02-15-04"
const secondKeyFormat = "2006-01-02-15-04-05"

var eventMap = make(map[string](map[string]int))

func GetEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := vars["event"]

	now := time.Now()
	hourKey := now.Format(hourKeyFormat)
	minuteKey := now.Format(minuteKeyFormat)
	secondKey := now.Format(secondKeyFormat)

	valuesToDate := map[string]int{
		"hour":   0,
		"minute": 0,
		"second": 0,
	}

	trailingCountPerSecond := 0
	trailingCountPerMinute := 0
	trailingCountPerHour := 0

	if _, eventExists := eventMap[event]; eventExists {
		if hourToDate, hourKeyExists := eventMap[event][hourKey]; hourKeyExists {
			valuesToDate["hour"] = hourToDate
		}
		if minuteToDate, minuteKeyExists := eventMap[event][minuteKey]; minuteKeyExists {
			valuesToDate["minute"] = minuteToDate
		}
		if secondToDate, secondKeyExists := eventMap[event][secondKey]; secondKeyExists {
			valuesToDate["second"] = secondToDate
		}

		// Seconds
		oneSecondAgo := now.Add(-1 * time.Second)
		secondKey = oneSecondAgo.Format(secondKeyFormat)
		if secondValue, secondKeyExists := eventMap[event][secondKey]; secondKeyExists {
			trailingCountPerSecond += secondValue
		}

		twoSecondsAgo := now.Add(-2 * time.Second)
		secondKey = twoSecondsAgo.Format(secondKeyFormat)
		if secondValue, secondKeyExists := eventMap[event][secondKey]; secondKeyExists {
			trailingCountPerSecond += secondValue
		}

		threeSecondsAgo := now.Add(-3 * time.Second)
		secondKey = threeSecondsAgo.Format(secondKeyFormat)
		if secondValue, secondKeyExists := eventMap[event][secondKey]; secondKeyExists {
			trailingCountPerSecond += secondValue
		}

		// Minutes
		oneMinuteAgo := now.Add(-1 * time.Minute)
		minuteKey = oneMinuteAgo.Format(minuteKeyFormat)
		if minuteValue, minuteKeyExists := eventMap[event][minuteKey]; minuteKeyExists {
			trailingCountPerMinute += minuteValue
		}

		twoMinutesAgo := now.Add(-2 * time.Minute)
		minuteKey = twoMinutesAgo.Format(minuteKeyFormat)
		if minuteValue, minuteKeyExists := eventMap[event][minuteKey]; minuteKeyExists {
			trailingCountPerMinute += minuteValue
		}

		threeMinutesAgo := now.Add(-3 * time.Minute)
		minuteKey = threeMinutesAgo.Format(minuteKeyFormat)
		if minuteValue, minuteKeyExists := eventMap[event][minuteKey]; minuteKeyExists {
			trailingCountPerMinute += minuteValue
		}

		// Hours
		oneHourAgo := now.Add(-1 * time.Hour)
		hourKey = oneHourAgo.Format(hourKeyFormat)
		if hourValue, hourKeyExists := eventMap[event][hourKey]; hourKeyExists {
			trailingCountPerHour += hourValue
		}

		twoHoursAgo := now.Add(-2 * time.Hour)
		hourKey = twoHoursAgo.Format(hourKeyFormat)
		if hourValue, hourKeyExists := eventMap[event][hourKey]; hourKeyExists {
			trailingCountPerHour += hourValue
		}

		threeHoursAgo := now.Add(-3 * time.Hour)
		hourKey = threeHoursAgo.Format(hourKeyFormat)
		if hourValue, hourKeyExists := eventMap[event][hourKey]; hourKeyExists {
			trailingCountPerHour += hourValue
		}
	}

	eventResponse := EventResponse{
		valuesToDate["second"],
		valuesToDate["minute"],
		valuesToDate["hour"],
		float32(trailingCountPerSecond) / float32(3.0),
		float32(trailingCountPerMinute) / float32(3.0),
		float32(trailingCountPerHour) / float32(3.0),
	}

	jsonData, err := json.Marshal(eventResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func updateCountForEvent(event string, t time.Time) {
	hourKey := t.Format(hourKeyFormat)
	minuteKey := t.Format(minuteKeyFormat)
	secondKey := t.Format(secondKeyFormat)

	dateMap, ok := eventMap[event]
	if ok == false {
		eventMap[event] = make(map[string]int)
		dateMap = eventMap[event]
	}

	if _, ok := dateMap[hourKey]; ok == false {
		dateMap[hourKey] = 0
	}
	dateMap[hourKey] += 1

	if _, ok := dateMap[minuteKey]; ok == false {
		dateMap[minuteKey] = 0
	}
	dateMap[minuteKey] += 1

	if _, ok := dateMap[secondKey]; ok == false {
		dateMap[secondKey] = 0
	}
	dateMap[secondKey] += 1

	// TODO replace eventMap with itself after filtering for relevant date keys
}

func PostEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := vars["event"]

	t := time.Now()
	updateCountForEvent(event, t)

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

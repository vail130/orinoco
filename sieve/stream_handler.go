package sieve

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"

	"github.com/vail130/orinoco/timeutils"
)

type Event struct {
	Stream string
	Data   []byte
}

type StreamSummary struct {
	Stream                   string  `json:"stream"`
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

// Reference time for formats: Mon Jan 2 15:04:05 -0700 MST 2006
const dateKeyFormat = "2006-01-02"
const hourKeyFormat = "2006-01-02-15"
const minuteKeyFormat = "2006-01-02-15-04"
const secondKeyFormat = "2006-01-02-15-04-05"

var dateMap = make(map[string](map[string](map[string]int)))
var dateKeyMap = make(map[string]time.Time)

func trackStreamForTime(stream string, t time.Time) {
	streamMap := GetStreamMapForTime(t)
	deleteObsoleteDateKeysForTime(t)

	dateMap, ok := streamMap[stream]
	if ok == false {
		streamMap[stream] = make(map[string]int)
		dateMap = streamMap[stream]
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

func GetTimestampForRequest(queryValues url.Values, data []byte) time.Time {
	if !sieveConfig.IsTestEnv {
		return timeutils.UtcNow()
	}

	if timestampString, ok := queryValues["timestamp"]; ok {
		if timestamp, err := time.Parse(time.RFC3339, timestampString[0]); err == nil {
			return timestamp.UTC()
		}
	}

	if len(data) > 0 {
		var f interface{}
		err := json.Unmarshal(data, &f)
		if err != nil {
			return timeutils.UtcNow()
		}

		if timestampString, ok := f.(map[string]interface{})["timestamp"]; ok {
			if timestamp, err := time.Parse(time.RFC3339, timestampString.(string)); err == nil {
				return timestamp
			}
		}
	}

	return timeutils.UtcNow()
}

func PostStreamHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		data = make([]byte, 0)
	}

	vars := mux.Vars(r)
	streamName := vars["stream"]

	t := GetTimestampForRequest(r.URL.Query(), data)
	trackStreamForTime(streamName, t)

	event := &Event{
		streamName,
		data,
	}

	sieveConfig.EventChannel <- event

	w.WriteHeader(http.StatusOK)
}
